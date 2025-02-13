package updatedpods

import (
	"context"

	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type InstrumentationConfigReconciler struct {
	client.Client
}

func (r *InstrumentationConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	pw, err := workload.PodWorkloadFromInstrumentationConfigName(req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, reconcile.TerminalError(err)
	}
	return reconcileWorkload(ctx, r.Client, pw)
}
