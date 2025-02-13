package workload

import (
	"github.com/odigos-io/odigos/api/k8sconsts"
	corev1 "k8s.io/api/core/v1"
)

func PodWorkloadFromPodObject(pod *corev1.Pod) (*k8sconsts.PodWorkload, error) {

	for _, owner := range pod.OwnerReferences {

		workloadName, workloadKind, err := GetWorkloadFromOwnerReference(owner)
		if IsErrorKindNotSupported(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		return &k8sconsts.PodWorkload{
			Name:      workloadName,
			Kind:      workloadKind,
			Namespace: pod.Namespace,
		}, nil
	}

	// Pod does not necessarily have to be managed by a controller
	return nil, ErrKindNotSupported
}
