package instrumentation_ebpf

import (
	"context"
	"errors"

	"github.com/odigos-io/odigos/common"

	odigosv1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	odigosK8s "github.com/odigos-io/odigos/k8sutils/pkg/container"
	"github.com/odigos-io/odigos/k8sutils/pkg/workload"
	"github.com/odigos-io/odigos/odiglet/pkg/ebpf"
	"github.com/odigos-io/odigos/odiglet/pkg/process"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func cleanupEbpf(directors ebpf.DirectorsMap, name types.NamespacedName) {
	// cleanup using all available directors
	// the Cleanup method is idempotent, so no harm in calling it multiple times
	for _, director := range directors {
		director.Cleanup(name)
	}
}

func instrumentPodWithEbpf(ctx context.Context, pod *corev1.Pod, directors ebpf.DirectorsMap, runtimeDetails *odigosv1.InstrumentedApplication, podWorkload *workload.PodWorkload) (error, bool) {
	logger := log.FromContext(ctx)
	podUid := string(pod.UID)
	instrumentedEbpf := false
	ignoredContainers := getIgnoredContainers(runtimeDetails)

	for _, container := range pod.Spec.Containers {
		// Ignored containers should not have a device but double check just in case
		if _, ignored := ignoredContainers[container.Name]; ignored {
			continue
		}

		var directorsToApply []ebpf.Director
		deviceName := odigosK8s.PodContainerDeviceName(container)
		if deviceName == nil {
			continue // the container does not have an instrumentation device
		}

		// hack to apply eBPF directors for all languages.
		// this is currently hard-coded which is not ideal.
		// we need to make it generic if we choose to keep it
		switch *deviceName {
		case common.OdigosInstrumentationPluginAllNativeCommunity:
			ebpfGo := directors[ebpf.DirectorKey{Language: common.GoProgrammingLanguage, OtelSdk: common.OtelSdkEbpfCommunity}]
			directorsToApply = []ebpf.Director{ebpfGo}
		case common.OdigosInstrumentationPluginAllNativeEnterprise:
			ebpfGo := directors[ebpf.DirectorKey{Language: common.GoProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			ebpfPython := directors[ebpf.DirectorKey{Language: common.PythonProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			ebpfNode := directors[ebpf.DirectorKey{Language: common.JavascriptProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			javaNativeInst := directors[ebpf.DirectorKey{Language: common.JavaProgrammingLanguage, OtelSdk: common.OtelSdkNativeEnterprise}]
			directorsToApply = []ebpf.Director{ebpfGo, ebpfPython, ebpfNode, javaNativeInst}
		case common.OdigosInstrumentationPluginAllEbpfEnterprise:
			ebpfGo := directors[ebpf.DirectorKey{Language: common.GoProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			ebpfPython := directors[ebpf.DirectorKey{Language: common.PythonProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			ebpfNode := directors[ebpf.DirectorKey{Language: common.JavascriptProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			javaEbpfInst := directors[ebpf.DirectorKey{Language: common.JavaProgrammingLanguage, OtelSdk: common.OtelSdkEbpfEnterprise}]
			directorsToApply = []ebpf.Director{ebpfGo, ebpfPython, ebpfNode, javaEbpfInst}
		default:
			language, sdk, found := odigosK8s.GetLanguageAndOtelSdk(container)
			if !found {
				continue
			}
			director := directors[ebpf.DirectorKey{Language: language, OtelSdk: sdk}]
			if director == nil {
				continue
			}
			directorsToApply = []ebpf.Director{director}
		}

		details, err := process.FindAllInContainer(podUid, container.Name)
		if err != nil {
			logger.Error(err, "error finding processes")
			return err, instrumentedEbpf
		}

		var errs []error
	processesLoop:
		for _, d := range details {
			podDetails := types.NamespacedName{
				Namespace: pod.Namespace,
				Name:      pod.Name,
			}

			for _, director := range directorsToApply {
				// Hack - Nginx instrumentation needs to be on Master Pid which is the min Pid in the container
				if !director.ShouldInstrument(d.ProcessID, details) {
					continue processesLoop
				}

				err = director.Instrument(ctx, d.ProcessID, podDetails, podWorkload, podWorkload.Name, container.Name)

				if err != nil {
					logger.Error(err, "error initiating process instrumentation", "pid", d.ProcessID)
					errs = append(errs, err)
					continue processesLoop
				}
				instrumentedEbpf = true
			}
		}

		// Failed to instrument all processes in the container
		if len(errs) > 0 && len(errs) == len(details) {
			return errors.Join(errs...), instrumentedEbpf
		}
	}
	return nil, instrumentedEbpf
}

func getIgnoredContainers(instApp *odigosv1.InstrumentedApplication) map[string]struct{} {
	ignoredContainers := make(map[string]struct{})
	for _, container := range instApp.Spec.RuntimeDetails {
		if container.Language == common.IgnoredProgrammingLanguage {
			ignoredContainers[container.ContainerName] = struct{}{}
		}
	}
	return ignoredContainers
}
