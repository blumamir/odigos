package actions

// Used to mark sources to avoid default templatization on error.
// Publicly accessible services are commonly being "tested" by malicious actors
// with irrelevant or garbage requests that can contaminate the url-templatization process
// leading to high-cardinality of templated routes.
// when no custom templatization rule matched, the request status code is checked, and skipped based on the config.
// either skipForNonSuccessCodes or statusCodes must be set for it to take effect.
// if both are set, skipForNonSuccessCodes takes precedence.
//
// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type DefaultTemplatizationSkipPolicyConfig struct {

	// If set to true, default templatization will be skipped for any non-success HTTP status code (i.e., any code not in the 2xx range).
	// When this is true, the StatusCodes list is ignored.
	SkipForNonSuccessCodes bool `json:"skipForNonSuccessCodes,omitempty"`

	// the http status codes for which the default templatization will be skipped.
	// for example: [404, 401].
	// advanced users can use to cherry pick specific codes for tailoring to specific use cases.
	SkipHttpStatusCodes []int `json:"skipHttpStatusCodes,omitempty"`
}

// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type DefaultTemplatizationConfig struct {

	// if set to true, the default templatization will be disabled for the services in the scope.
	// in case of conflict (if other action set this to false), the default templatization will be disabled.
	Disabled bool `json:"disabled,omitempty"`

	// config if default templatization should be skipped on error.
	// use it when sources scope describes a service that is publicly accessible to the internet
	// to filter garbage requests that can contaminate the url-templatization process.
	SkipPolicy *DefaultTemplatizationSkipPolicyConfig `json:"skipPolicy,omitempty"`
}

// +kubebuilder:object:generate=true
type UrlTemplatizationConfig struct {
	// Template rules to apply to URLs
	Templates []string `json:"templatizationRules,omitempty"`

	// configurations for default templatization.
	// default templatization is applied on a single http span if none of the custom templatization rules matched.
	Default *DefaultTemplatizationConfig `json:"default,omitempty"`
}

// URLTemplatizationRuleCertainty describes how confident live traffic learning is in a templated path segment.
// +kubebuilder:validation:Enum=High;Moderate
type URLTemplatizationRuleCertainty string

const (
	URLTemplatizationRuleCertaintyHigh     URLTemplatizationRuleCertainty = "High"
	URLTemplatizationRuleCertaintyModerate URLTemplatizationRuleCertainty = "Moderate"
)

// URLTemplatizationSegment describes one templated path segment in a learned rule.
//
// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type URLTemplatizationSegment struct {
	// TemplateName is the name inside braces in the template rule (e.g. "id" for "{id}").
	// When the same name appears more than once in a template, entries are ordered by
	// appearance of templated segments in the rule (left to right).
	// For a wildcard segment the value is "*".
	TemplateName string `json:"templateName"`

	// Certainty indicates how confident live traffic learning is in this templated segment.
	Certainty URLTemplatizationRuleCertainty `json:"certainty"`

	// Examples are concrete path-segment values observed for this templated segment.
	// Learning keeps up to 5 examples.
	Examples []string `json:"examples,omitempty"`
}

// URLTemplatizationLearnedRule is one learned URL templatization recommendation.
// Presence in status means the recommendation is pending; accepting or rejecting removes it.
//
// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type URLTemplatizationLearnedRule struct {
	// Template is the recommended URL template rule (e.g. "/users/{id}/orders/{orderId}").
	Template string `json:"template"`

	// Reason is a human-readable explanation of why this rule was suggested.
	Reason string `json:"reason,omitempty"`

	// Segments lists the templated segments in this rule.
	Segments []URLTemplatizationSegment `json:"segments,omitempty"`
}

// URLTemplatizationContainerFindings holds learned URL templatization recommendations for one container.
//
// +kubebuilder:object:generate=true
// +kubebuilder:deepcopy-gen=true
type URLTemplatizationContainerFindings struct {
	// ContainerName is the name of the container within the workload.
	ContainerName string `json:"containerName"`

	// Rules are the learned URL templatization recommendations for this container.
	Rules []URLTemplatizationLearnedRule `json:"rules,omitempty"`
}
