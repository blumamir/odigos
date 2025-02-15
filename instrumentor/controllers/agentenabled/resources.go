package agentenabled

import (
	"context"
	"slices"

	"github.com/odigos-io/odigos/api/k8sconsts"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/instrumentor/controllers/utils"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"
	k8sutils "github.com/odigos-io/odigos/k8sutils/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// fetches the relevant resources for reconciliation of the current workload object.
// if err is returned, the reconciliation should be retried.
// the function can return without error, but some of the resources might be nil.
// this indicates that they were not found, but it is valid.
func getRelevantResources(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) (
	*odigosv1.CollectorsGroup,
	*[]odigosv1.InstrumentationRule,
	*common.OdigosConfiguration,
	error) {

	cg, err := getCollectorsGroup(ctx, c)
	if err != nil {
		return nil, nil, nil, err
	}

	irls, err := getRelevantInstrumentationRules(ctx, c, pw)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: we are yaml unmarshalling the configmap data for every workload, this is not efficient
	// can we cache the configmap data in the controller?
	effectiveConfig, err := k8sutils.GetCurrentOdigosConfig(ctx, c)
	if err != nil {
		return nil, nil, nil, err
	}

	return cg, irls, &effectiveConfig, nil
}

func getCollectorsGroup(ctx context.Context, c client.Client) (*odigosv1.CollectorsGroup, error) {
	cg := odigosv1.CollectorsGroup{}
	err := c.Get(ctx, client.ObjectKey{Namespace: env.GetCurrentNamespace(), Name: k8sconsts.OdigosNodeCollectorCollectorGroupName}, &cg)
	if err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return &cg, nil
}

func getRelevantInstrumentationRules(ctx context.Context, c client.Client, pw k8sconsts.PodWorkload) (*[]odigosv1.InstrumentationRule, error) {
	relevantIr := []odigosv1.InstrumentationRule{}
	irList := odigosv1.InstrumentationRuleList{}
	err := c.List(ctx, &irList)
	if err != nil {
		return nil, err
	}

	for i := range irList.Items {
		ir := &irList.Items[i]

		if !utils.IsWorkloadParticipatingInRule(pw, ir) {
			continue
		}

		if ir.Spec.OtelSdks == nil && ir.Spec.OtelDistros == nil {
			// we only care about otel sdks rules at the moment.
			// no need to process other rules.
			continue
		}

		relevantIr = append(relevantIr, *ir)
	}

	// sort the instrumentation rules to apply them at deterministic order
	// and to have more relevant rules applied last.
	// this insures rules from profiles are applied first, so they can be overridden by more specific rules.
	// TODO: this assumes we have very few rules (expecting few rules at most and no more than 10).
	// if in the future we have many rules, we should consider a more efficient sorting algorithm.
	slices.SortFunc(relevantIr, func(ir1, ir2 odigosv1.InstrumentationRule) int {

		mb1 := ir1.Labels[k8sconsts.OdigosProfilesManagedByLabel]
		mb2 := ir2.Labels[k8sconsts.OdigosProfilesManagedByLabel]

		ir1FromProfile := mb1 == k8sconsts.OdigosProfilesManagedByValue
		ir2FromProfile := mb2 == k8sconsts.OdigosProfilesManagedByValue

		if ir1FromProfile && !ir2FromProfile {
			return -1
		} else if !ir1FromProfile && ir2FromProfile {
			return 1
		} else {
			// if both are from profile or both are not from profile, sort by creation timestamp
			return ir1.CreationTimestamp.Compare(ir2.CreationTimestamp.Time)
		}
	})

	return &relevantIr, nil
}
