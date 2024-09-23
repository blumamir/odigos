package instrumentationrules

// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type IgnoredContainers struct {
	ContainerNames []string `json:"containerNames"`
}
