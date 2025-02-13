package updatedpods

import (
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func SetupWithManager(mgr ctrl.Manager) error {
	err := builder.
		ControllerManagedBy(mgr).
		Named("updatedpods-pods").
		For(&corev1.Pod{}).
		Complete(&PodReconciler{
			Client: mgr.GetClient(),
		})
	if err != nil {
		return err
	}

	err = builder.
		ControllerManagedBy(mgr).
		Named("updatedpods-instrumentationconfig").
		For(&odigosv1.InstrumentationConfig{}).
		Complete(&InstrumentationConfigReconciler{
			Client: mgr.GetClient(),
		})
	if err != nil {
		return err
	}

	return nil
}
