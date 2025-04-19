package logs

import (
	actionv1 "github.com/odigos-io/odigos/api/actions/v1alpha1"
	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	odigospredicate "github.com/odigos-io/odigos/k8sutils/pkg/predicate"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func SetupWithManager(mgr ctrl.Manager) error {

	err := builder.
		ControllerManagedBy(mgr).
		Named("logs-collectorsgroup").
		For(&odigosv1.CollectorsGroup{}).
		Owns(&actionv1.K8sAttributesResolver{}). // reconcile the object if drifted
		WithEventFilter(&odigospredicate.OdigosCollectorsGroupClusterPredicate).
		Complete(&CollectorsGroupReconciler{
			Client: mgr.GetClient(),
		})
	if err != nil {
		return err
	}

	return nil
}
