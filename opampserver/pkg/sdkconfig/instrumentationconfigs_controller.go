package sdkconfig

import (
	"context"
	"encoding/json"

	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	"github.com/odigos-io/odigos/opampserver/pkg/connection"
	"github.com/odigos-io/odigos/opampserver/pkg/sdkconfig/configsections"
	"github.com/odigos-io/odigos/opampserver/protobufs"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstrumentationConfigReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	ConnectionCache *connection.ConnectionsCache
}

func (i *InstrumentationConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instrumentationConfig := &odigosv1.InstrumentationConfig{}
	err := i.Get(ctx, req.NamespacedName, instrumentationConfig)

	if err != nil {
		if apierrors.IsNotFound(err) {
			instrumentationConfig = nil
		} else {
			return ctrl.Result{}, err
		}
	}

	workloadName, workloadKind, err := workload.GetWorkloadInfoRuntimeName(req.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	podWorkload := common.PodWorkload{
		Namespace: req.Namespace,
		Kind:      workloadKind,
		Name:      workloadName,
	}

	instrumentationLibrariesConfig, err := configsections.CalcInstrumentationLibrariesRemoteConfig(ctx, i.Client, req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}
	instrumentationLibrariesConfigBytes, err := json.Marshal(instrumentationLibrariesConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	updatedConfigMapEntries := protobufs.AgentConfigMap{
		ConfigMap: map[string]*protobufs.AgentConfigFile{
			string(configsections.RemoteConfigInstrumentationLibrariesConfigSectionName): {
				Body:        instrumentationLibrariesConfigBytes,
				ContentType: "application/json",
			},
		},
	}

	i.ConnectionCache.UpdateWorkloadRemoteConfigByKeys(podWorkload, &updatedConfigMapEntries)

	return ctrl.Result{}, nil
}

func (i *InstrumentationConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odigosv1.InstrumentationConfig{}).
		Complete(i)
}
