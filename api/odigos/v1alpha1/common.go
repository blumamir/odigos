package v1alpha1

// Attribute is a key-value pair that conforms to the OpenTelemetry attributes specification
// and semantic conventions.
// while opentelemetry attributes supports a wide range of types for attribute value, this struct only supports strings.
type Attribute struct {
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`
	// +required
	// +kubebuilder:validation:Required
	Value string `json:"value"`
}
