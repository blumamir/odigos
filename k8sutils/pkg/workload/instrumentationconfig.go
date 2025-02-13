package workload

import (
	"errors"
	"strings"

	"github.com/odigos-io/odigos/api/k8sconsts"
)

func CalculateWorkloadRuntimeObjectName[T string | k8sconsts.WorkloadKind | k8sconsts.WorkloadKindLowerCase](
	workloadName string, workloadKind T) string {
	return strings.ToLower(string(workloadKind) + "-" + workloadName)
}

func InstrumentationConfigNameFromPodWorkload(pw k8sconsts.PodWorkload) string {
	return CalculateWorkloadRuntimeObjectName(pw.Name, pw.Kind)
}

func ExtractWorkloadInfoFromRuntimeObjectName(runtimeObjectName string) (workloadName string, workloadKind k8sconsts.WorkloadKind, err error) {
	parts := strings.SplitN(runtimeObjectName, "-", 2)
	if len(parts) != 2 {
		err = errors.New("invalid workload runtime object name, missing hyphen")
		return
	}

	// convert the lowercase kind to pascal case and validate it
	workloadKindLowerCase := k8sconsts.WorkloadKindLowerCase(parts[0])
	workloadKind = WorkloadKindFromLowerCase(workloadKindLowerCase)
	if workloadKind == "" {
		err = ErrKindNotSupported
		return
	}

	workloadName = parts[1]

	return
}

func PodWorkloadFromInstrumentationConfigName(icName string, ns string) (k8sconsts.PodWorkload, error) {
	workloadName, workloadKind, err := ExtractWorkloadInfoFromRuntimeObjectName(icName)
	if err != nil {
		return k8sconsts.PodWorkload{}, err
	}

	pw := k8sconsts.PodWorkload{
		Kind:      workloadKind,
		Name:      workloadName,
		Namespace: ns,
	}

	return pw, nil
}
