package sdkconfig

import (
	"context"

	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"github.com/odigos-io/odigos/common"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InstrumentedApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *InstrumentedApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	// get the language from the instrumented application
	instrumentedApplication := &odigosv1.InstrumentedApplication{}
	err := r.Get(ctx, req.NamespacedName, instrumentedApplication)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	supportedLanguages := make([]common.ProgrammingLanguage, 0, 1)
	for _, containerRuntimeDetails := range instrumentedApplication.Spec.RuntimeDetails {
		language := containerRuntimeDetails.Language
		if language == common.JavascriptProgrammingLanguage {
			supportedLanguages = append(supportedLanguages, language)
		}
	}

	// extract the workload info from the instrumented application name
	// and fetch the workload object
	workloadObj := &client.Object{}
	err = getWorkloadObject(ctx, r.Client, workloadObj)

	sdkConfig, err := calculateInstrumentationConfig(ctx, supportedLanguages)

	return ctrl.Result{}, nil
}

func getWorkloadObject(ctx context.Context, k8sClient client.Client, req ctrl.Request, obj client.Object) error {
	return k8sClient.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, obj)
}

func calculateInstrumentationConfig(ctx context.Context, supportedLanguages []common.ProgrammingLanguage) (odigosv1.InstrumentationConfigSpec, error) {

	// collect all instrumentation config from the various sources
	// and build one single instrumentation config object with everything related to this workload

	// fetch the workload object to check for service name configuration

	sdkConfigs := make([]odigosv1.SdkConfig, 0, len(supportedLanguages))
	for _, language := range supportedLanguages {
		sdkConfig := odigosv1.SdkConfig{
			Language: language,
		}
		sdkConfigs = append(sdkConfigs, sdkConfig)
	}

	instrumentationConfigSpec := odigosv1.InstrumentationConfigSpec{
		// RuntimeDetailsInvalidated we do not populate or control this field here.
		// at the time of writing, it is handled by the startlangdetection controllers

		// Config is not populate it by this controller.
		// TODO: populate it here, or migrate it into the sdk section

		SdkConfigs: sdkConfigs,
	}

	return instrumentationConfigSpec, nil
}
