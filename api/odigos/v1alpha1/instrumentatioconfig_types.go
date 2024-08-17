package v1alpha1

import (
	"github.com/odigos-io/odigos/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// InstrumentationConfig is the Schema for the instrumentationconfig API
type InstrumentationConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec InstrumentationConfigSpec `json:"spec,omitempty"`
}

// Config for the OpenTelemeetry SDKs that should be applied to a workload.
// The workload is identified by the owner reference
type InstrumentationConfigSpec struct {
	// true when the runtime details are invalidated and should be recalculated
	RuntimeDetailsInvalidated bool `json:"runtimeDetailsInvalidated,omitempty"`

	// config for this workload.
	// the config is a list to allow for multiple config options and values to be applied.
	// the list is processed in order, and the first matching config is applied.
	Config []WorkloadInstrumentationConfig `json:"config,omitempty"`

	// Configuration for the OpenTelemetry SDKs that this workload should use.
	// The SDKs are identified by the programming language they are written in.
	// TODO: consider adding more granular control over the SDKs, such as community/enterprise, native/ebpf.
	SdkConfigs []SdkConfig `json:"sdkConfigs,omitempty"`
}

type SdkConfig struct {

	// The language of the SDK being configured
	Language common.ProgrammingLanguage `json:"language"`

	// configurations for the instrumentation libraries the the SDK should use
	InstrumentationLibraryConfigs []InstrumentationLibraryConfig `json:"instrumentationLibraryConfigs"`

	// ResourceAttributes is a list of attributes that should be added to the telemetry
	// emitted from this SDK. These attributes are added to all spans, along with any other attributes
	// that are detected regardless, like process.*, k8s.*, etc.
	// if you want to control the `service.name` resource attribute value, pass it here.
	ResourceAttributes []Attribute `json:"resourceAttributes"`

	// general configuration for recording traces in the SDK level.
	Traces TraceSignalGeneralConfig `json:"traces"`
}

type InstrumentationLibraryConfig struct {
	// The name of the instrumentation library
	// - Node.js: The name of the npm package: `@opentelemetry/instrumentation-<name>`
	InstrumentationLibraryName string `json:"instrumentationLibraryName"`

	TraceConfig *InstrumentationLibraryConfigTraces `json:"traceConfig,omitempty"`
}

type InstrumentationLibraryConfigTraces struct {
	// Whether the instrumentation library is enabled to record traces.
	// When false, it is expected that the instrumentation library does not produce any spans regardless of any other configuration.
	// When true, the instrumentation library should produce spans according to the other configuration options.
	// If not specified, the default value for this signal should be used (whether to enable libraries by default or not).
	Enabled *bool `json:"enabled,omitempty"`
}

// general configuration in the SDK level for trace signal.
type TraceSignalGeneralConfig struct {

	// general flag to control if the sdk should produce traces or not.
	// this can be used to globally disabled traces when the receiver is not setup to receive traces,
	// if there is not destination for the traces.
	Enabled bool `json:"enabled"`

	// by using this value, one can choose the behavior for instrumentation libraries
	// for which there is no explicit enabled configuration.
	// simplest (black list) - one can set this value to true to allow all instrumentation libraries to produce traces, unless explicitly disabled.
	// safest (white list)- by default all instrumentation libraries are disabled, unless explicitly enabled.
	DefaultEnabledValue bool `json:"defaultEnabledValue"`
}

// WorkloadInstrumentationConfig defined a single config option to apply
// on a workload, along with it's value, filters and instrumentation libraries
type WorkloadInstrumentationConfig struct {

	// OptionKey is the name of the option
	// This value is transparent to the CRD and is passed as-is to the SDK.
	OptionKey string `json:"optionKey"`

	// This option allow to specify the config option for a specific span kind
	// for example, only to client spans or only to server spans.
	// it the span kind is not specified, the option will apply to all spans.
	SpanKind common.SpanKind `json:"spanKind,omitempty"`

	// OptionValueBoolean is the boolean value of the option if it is a boolean
	OptionValueBoolean bool `json:"optionValueBoolean,omitempty"`

	// a list of instrumentation libraries to apply this setting to
	// if a library is not in this list, the setting should not apply to it
	// and should be cleared.
	InstrumentationLibraries []InstrumentationLibrary `json:"instrumentationLibraries"`
}

// InstrumentationLibrary represents a library for instrumentation
type InstrumentationLibrary struct {
	// Language is the programming language of the library
	Language common.ProgrammingLanguage `json:"language"`

	// InstrumentationLibraryName is the name of the instrumentation library
	InstrumentationLibraryName string `json:"instrumentationLibraryName"`
}

// +kubebuilder:object:root=true

// InstrumentationConfigList contains a list of InstrumentationOption
type InstrumentationConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstrumentationConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InstrumentationConfig{}, &InstrumentationConfigList{})
}
