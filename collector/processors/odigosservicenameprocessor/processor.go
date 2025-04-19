package odigosservicenameprocessor

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	semconv "go.opentelemetry.io/collector/semconv/v1.27.0"
	"go.uber.org/zap"

	"github.com/odigos-io/odigos/api/k8sconsts"
)

type serviceNameProcessor struct {
	logger *zap.Logger
}

func newServiceNameProcessor(set processor.Settings, config *Config) (*serviceNameProcessor, error) {
	return &serviceNameProcessor{
		logger: set.Logger,
	}, nil
}

func (p *serviceNameProcessor) processLogs(ctx context.Context, ls plog.Logs) (plog.Logs, error) {
	rl := ls.ResourceLogs()
	for i := 0; i < rl.Len(); i++ {
		err := p.processResource(ctx, rl.At(i).Resource())
		if err != nil {
			return ls, err
		}
	}
	return ls, nil
}

// in order to resolve odigos "service.name", we need to first get the workload info:
// namespace, kind and name.
// we search for them on the resource attributes, and return a PodWorkload object or nil if not found.
func resourceAttributesToPodWorkload(resourceAttributes pcommon.Map) *k8sconsts.PodWorkload {
	namespace, ok := resourceAttributes.Get(semconv.AttributeK8SNamespaceName)
	if !ok {
		return nil
	}

	deploymentName, ok := resourceAttributes.Get(semconv.AttributeK8SDeploymentName)
	if ok {
		return &k8sconsts.PodWorkload{
			Namespace: namespace.Str(),
			Name:      deploymentName.Str(),
			Kind:      k8sconsts.WorkloadKindDeployment,
		}
	}

	statefulSetName, ok := resourceAttributes.Get(semconv.AttributeK8SStatefulSetName)
	if ok {
		return &k8sconsts.PodWorkload{
			Namespace: namespace.Str(),
			Name:      statefulSetName.Str(),
			Kind:      k8sconsts.WorkloadKindStatefulSet,
		}
	}

	daemonSetName, ok := resourceAttributes.Get(semconv.AttributeK8SDaemonSetName)
	if ok {
		return &k8sconsts.PodWorkload{
			Namespace: namespace.Str(),
			Name:      daemonSetName.Str(),
			Kind:      k8sconsts.WorkloadKindDaemonSet,
		}
	}

	// if neither of the above attributes are found, we return nil
	// this happens when the workload name is not set on the resource attributes
	return nil
}

func (p *serviceNameProcessor) processResource(ctx context.Context, res pcommon.Resource) error {
	podWorkload := resourceAttributesToPodWorkload(res.Attributes())
	if podWorkload == nil {
		// if we don't have a pod workload, we can't set the service name
		// we log a warning and return
		p.logger.Warn("no pod workload found, skipping service name processor")
		return nil
	}
	// test: add service name with static value
	res.Attributes().PutStr(semconv.AttributeServiceName, podWorkload.Name)
	return nil
}
