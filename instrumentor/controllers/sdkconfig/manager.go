package sdkconfig

import (
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func SetupWithManager(mgr ctrl.Manager) error {

	err := builder.
		ControllerManagedBy(mgr).
		For(&odigosv1.InstrumentedApplication{}).
		Complete(&InstrumentedApplicationReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
	if err != nil {
		return err
	}

	// err = builder.
	// 	ControllerManagedBy(mgr).
	// 	For(&appsv1.Deployment{}).
	// 	Complete(&DeploymentReconciler{
	// 		Client: mgr.GetClient(),
	// 		Scheme: mgr.GetScheme(),
	// 	})
	// if err != nil {
	// 	return err
	// }

	// err = builder.
	// 	ControllerManagedBy(mgr).
	// 	For(&appsv1.DaemonSet{}).
	// 	Complete(&DaemonSetReconciler{
	// 		Client: mgr.GetClient(),
	// 		Scheme: mgr.GetScheme(),
	// 	})
	// if err != nil {
	// 	return err
	// }

	// err = builder.
	// 	ControllerManagedBy(mgr).
	// 	For(&appsv1.StatefulSet{}).
	// 	Complete(&StatefulSetReconciler{
	// 		Client: mgr.GetClient(),
	// 		Scheme: mgr.GetScheme(),
	// 	})
	// if err != nil {
	// 	return err
	// }

	return nil
}
