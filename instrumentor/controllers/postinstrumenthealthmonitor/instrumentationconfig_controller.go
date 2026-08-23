package postinstrumenthealthmonitor

import (
	"context"

	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstrumentationConfigController struct {
	client.Client
	monitorEnabled *monitorEnabledState
}

func (r *InstrumentationConfigController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if !r.monitorEnabled.enabled.Load() {
		return ctrl.Result{}, nil
	}

	pw, err := workload.ExtractWorkloadInfoFromRuntimeObjectName(req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = syncWorkload(ctx, r.Client, pw)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
