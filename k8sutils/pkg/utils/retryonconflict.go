package utils

import (
	nativeErrors "errors"

	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// K8SUpdateErrorHandler is a helper function to handle k8s update errors.
// It returns a reconcile.Result and an error
// If the error is a conflict error, it returns a requeue result without an error
// If the error is a not found error, it returns an empty result without an error
// For other errors, it returns the error as is - which will cause a requeue
func K8SUpdateErrorHandler(err error) (reconcile.Result, error) {
	if errors.IsConflict(err) {
		// For conflict errors, requeue without returning an error.
		// this is so that we don't have errors and stack traces in the logs for valid scenario
		return reconcile.Result{Requeue: true}, nil
	}
	if errors.IsNotFound(err) {
		// For not found errors, ignore
		return reconcile.Result{}, nil
	}

	if nativeErrors.Is(err, ErrOtherAgentRun) {
		// For other agent run no need to log the stack trace
		return reconcile.Result{}, nil
	}
	// For other errors, return as is (will log the stack trace)
	return reconcile.Result{}, err
}

// MergeReconcileResults combines res into aggregated.
func MergeReconcileResults(aggregated, res reconcile.Result) reconcile.Result {
	if res.Requeue {
		aggregated.Requeue = true
		return aggregated
	}
	if res.RequeueAfter == 0 {
		return aggregated
	}
	if aggregated.RequeueAfter == 0 || res.RequeueAfter < aggregated.RequeueAfter {
		aggregated.RequeueAfter = res.RequeueAfter
	}
	return aggregated
}

// AggregateK8SUpdateError applies K8SUpdateErrorHandler to err and merges the result
// into aggregated reconcile state. Use in loops over many updates instead of joining
// errors and calling K8SUpdateErrorHandler once at the end.
func AggregateK8SUpdateError(aggregatedResult reconcile.Result, aggregatedErr error, err error) (reconcile.Result, error) {
	if err == nil {
		return aggregatedResult, aggregatedErr
	}
	res, handledErr := K8SUpdateErrorHandler(err)
	aggregatedResult = MergeReconcileResults(aggregatedResult, res)
	if handledErr != nil {
		aggregatedErr = nativeErrors.Join(aggregatedErr, handledErr)
	}
	return aggregatedResult, aggregatedErr
}
