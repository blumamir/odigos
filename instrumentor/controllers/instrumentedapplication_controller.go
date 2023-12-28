/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	odigosv1 "github.com/keyval-dev/odigos/api/odigos/v1alpha1"
	"github.com/keyval-dev/odigos/common/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// InstrumentedApplicationReconciler reconciles a InstrumentedApplication object
type InstrumentedApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=odigos.io,resources=instrumentedapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=odigos.io,resources=instrumentedapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=odigos.io,resources=instrumentedapplications/finalizers,verbs=update
//+kubebuilder:rbac:groups=odigos.io,resources=odigosconfigurations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/status,verbs=get;update;patch

// Reconcile is responsible for instrumenting deployment/statefulset/daemonset. In order for instrumentation to happen two things must be true:
// 1. InstrumentedApplication must have at least one language specified
// 2. Data collection pods must be running (DataCollection CollectorsGroup .status.ready == true)
func (r *InstrumentedApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	err := reconcileSingleInstrumentedApplication(ctx, r.Client, logger, req.NamespacedName)
	return ctrl.Result{}, err
}

// this function is extracted so we can call it from other reconcilers like when odigos config changes
func reconcileSingleInstrumentedApplication(ctx context.Context, kubeClient client.Client, logger logr.Logger, instrumentedApplicationName types.NamespacedName) error {
	var runtimeDetails odigosv1.InstrumentedApplication
	err := kubeClient.Get(ctx, instrumentedApplicationName, &runtimeDetails)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "error fetching instrumented application")
			return err
		}

		// runtime details deleted: remove instrumentation from resource requests
		err = removeInstrumentation(logger, ctx, kubeClient, instrumentedApplicationName)
		if err != nil {
			logger.Error(err, "error removing instrumentation")
			return err
		}

		return nil
	}

	if len(runtimeDetails.Spec.Languages) == 0 {
		err = removeInstrumentation(logger, ctx, kubeClient, instrumentedApplicationName)
		if err != nil {
			logger.Error(err, "error removing instrumentation")
			return err
		}

		return nil
	}

	if !isDataCollectionReady(ctx, kubeClient) {
		err := removeInstrumentation(logger, ctx, kubeClient, instrumentedApplicationName)
		if err != nil {
			logger.Error(err, "error removing instrumentation")
			return err
		}
	} else {
		err := instrument(logger, ctx, kubeClient, &runtimeDetails)
		if err != nil {
			logger.Error(err, "error instrumenting")
			return err
		}
	}

	return nil
}

func removeInstrumentation(logger logr.Logger, ctx context.Context, kubeClient client.Client, instrumentedApplicationName types.NamespacedName) error {
	name, kind, err := utils.GetTargetFromRuntimeName(instrumentedApplicationName.Name)
	if err != nil {
		return err
	}

	err = uninstrument(logger, ctx, kubeClient, instrumentedApplicationName.Namespace, name, kind)
	if err != nil {
		logger.Error(err, "error removing instrumentation")
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *InstrumentedApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&odigosv1.InstrumentedApplication{}).
		Complete(r)
}
