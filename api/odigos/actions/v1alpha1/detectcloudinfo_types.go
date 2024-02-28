package v1alpha1

import (
	"github.com/keyval-dev/odigos/common"
)

type AwsEC2Attributes struct {
	CloudProvider         bool `json:"cloudProvider,omitempty"`
	CloudPlatform         bool `json:"cloudPlatform,omitempty"`
	CloudAccountId        bool `json:"cloudAccountId,omitempty"`
	CloudRegion           bool `json:"cloudRegion,omitempty"`
	CloudAvailabilityZone bool `json:"cloudAvailabilityZone,omitempty"`
	HostId                bool `json:"hostId,omitempty"`
	HostImageId           bool `json:"hostImageId,omitempty"`
	HostName              bool `json:"hostName,omitempty"`
	HostType              bool `json:"hostType,omitempty"`
}

type DetectCloudInfoSpec struct {
	ActionName string                       `json:"actionName,omitempty"`
	Notes      string                       `json:"notes,omitempty"`
	Disabled   bool                         `json:"disabled,omitempty"`
	Signals    []common.ObservabilitySignal `json:"signals"`

	AwsEC2Attributes *AwsEC2Attributes `json:"awsEC2Attributes,omitempty"`
}
