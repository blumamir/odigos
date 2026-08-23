package postinstrumenthealthmonitor

import (
	"context"
	"errors"

	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/k8sutils/pkg/utils"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EffectiveConfigReconciler reacts to auto-rollback settings in the effective config
// and re-evaluates post-instrument health monitoring for all InstrumentationConfigs.
type EffectiveConfigReconciler struct {
	client.Client
	monitorEnabled *monitorEnabledState
}

func (r *EffectiveConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	effectiveConfig, err := utils.GetCurrentOdigosConfiguration(ctx, r.Client)
	if err != nil {
		return utils.K8SNoEffectiveConfigErrorHandler(err)
	}

	rollbackConfig, err := getAutoRollbackConfig(&effectiveConfig)
	if err != nil {
		return ctrl.Result{}, err
	}
	r.monitorEnabled.enabled.Store(!rollbackConfig.disabled)

	allInstrumentationConfigs := odigosv1.InstrumentationConfigList{}
	err = r.Client.List(ctx, &allInstrumentationConfigs)
	if err != nil {
		return ctrl.Result{}, err
	}

	aggregatedResult := ctrl.Result{}
	var allErrs error
	for _, ic := range allInstrumentationConfigs.Items {
		pw, err := workload.ExtractWorkloadInfoFromRuntimeObjectName(ic.Name, ic.Namespace)
		if err != nil {
			allErrs = errors.Join(allErrs, err)
			continue
		}

		_, err = syncWorkload(ctx, r.Client, pw)
		aggregatedResult, allErrs = utils.AggregateK8SUpdateError(aggregatedResult, allErrs, err)
	}

	return aggregatedResult, allErrs
}
