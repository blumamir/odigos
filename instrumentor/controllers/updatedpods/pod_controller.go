package updatedpods

import (
	"context"

	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type PodReconciler struct {
	client.Client
}

func (r *PodReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	var pod corev1.Pod
	err := r.Client.Get(ctx, req.NamespacedName, &pod)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	podWorkload, err := workload.PodWorkloadFromPodObject(&pod)
	if err != nil || podWorkload == nil {
		return reconcile.Result{}, reconcile.TerminalError(workload.IgnoreErrorKindNotSupported(err))
	}

	return reconcileWorkload(ctx, r.Client, *podWorkload)
}
