package instrumentationconfig

import (
	"context"

	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/k8sutils/pkg/consts"
	"github.com/odigos-io/odigos/k8sutils/pkg/env"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetEnabledSignals(ctx context.Context, k8sclient client.Client) ([]common.ObservabilitySignal, error) {
	nodeCollectorsGroup := &odigosv1.CollectorsGroup{}
	err := k8sclient.Get(ctx, client.ObjectKey{Namespace: env.GetCurrentNamespace(), Name: consts.OdigosNodeCollectorCollectorGroupName}, nodeCollectorsGroup)
	if err != nil {
		return nil, err
	}

	// At the moment, we use a homogeneous set of signals for all the applications.
	// thus - the workloads signals are the ones that node collectors group is set to collect.
	// in the future, we might want to have a more fine-grained control over the signals and which workload enables which signal
	enabledSignals := nodeCollectorsGroup.Status.ReceiverSignals
	return enabledSignals, nil
}
