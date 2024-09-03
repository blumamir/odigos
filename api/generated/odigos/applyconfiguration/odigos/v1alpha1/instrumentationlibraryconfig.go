/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	rulesv1alpha1 "github.com/odigos-io/odigos/api/rules/v1alpha1"
)

// InstrumentationLibraryConfigApplyConfiguration represents a declarative configuration of the InstrumentationLibraryConfig type for use
// with apply.
type InstrumentationLibraryConfigApplyConfiguration struct {
	InstrumentationLibraryId      *InstrumentationLibraryIdApplyConfiguration           `json:"libraryId,omitempty"`
	TraceConfig                   *InstrumentationLibraryConfigTracesApplyConfiguration `json:"traceConfig,omitempty"`
	HttpRequestPayloadCollection  *rulesv1alpha1.HttpPayloadCollectionRule              `json:"httpRequestPayloadCollection,omitempty"`
	HttpResponsePayloadCollection *rulesv1alpha1.HttpPayloadCollectionRule              `json:"httpResponsePayloadCollection,omitempty"`
	DbStatementPayloadCollection  *rulesv1alpha1.DbStatementPayloadCollectionRule       `json:"dbStatementPayloadCollection,omitempty"`
}

// InstrumentationLibraryConfigApplyConfiguration constructs a declarative configuration of the InstrumentationLibraryConfig type for use with
// apply.
func InstrumentationLibraryConfig() *InstrumentationLibraryConfigApplyConfiguration {
	return &InstrumentationLibraryConfigApplyConfiguration{}
}

// WithInstrumentationLibraryId sets the InstrumentationLibraryId field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the InstrumentationLibraryId field is set to the value of the last call.
func (b *InstrumentationLibraryConfigApplyConfiguration) WithInstrumentationLibraryId(value *InstrumentationLibraryIdApplyConfiguration) *InstrumentationLibraryConfigApplyConfiguration {
	b.InstrumentationLibraryId = value
	return b
}

// WithTraceConfig sets the TraceConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TraceConfig field is set to the value of the last call.
func (b *InstrumentationLibraryConfigApplyConfiguration) WithTraceConfig(value *InstrumentationLibraryConfigTracesApplyConfiguration) *InstrumentationLibraryConfigApplyConfiguration {
	b.TraceConfig = value
	return b
}

// WithHttpRequestPayloadCollection sets the HttpRequestPayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HttpRequestPayloadCollection field is set to the value of the last call.
func (b *InstrumentationLibraryConfigApplyConfiguration) WithHttpRequestPayloadCollection(value rulesv1alpha1.HttpPayloadCollectionRule) *InstrumentationLibraryConfigApplyConfiguration {
	b.HttpRequestPayloadCollection = &value
	return b
}

// WithHttpResponsePayloadCollection sets the HttpResponsePayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HttpResponsePayloadCollection field is set to the value of the last call.
func (b *InstrumentationLibraryConfigApplyConfiguration) WithHttpResponsePayloadCollection(value rulesv1alpha1.HttpPayloadCollectionRule) *InstrumentationLibraryConfigApplyConfiguration {
	b.HttpResponsePayloadCollection = &value
	return b
}

// WithDbStatementPayloadCollection sets the DbStatementPayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DbStatementPayloadCollection field is set to the value of the last call.
func (b *InstrumentationLibraryConfigApplyConfiguration) WithDbStatementPayloadCollection(value rulesv1alpha1.DbStatementPayloadCollectionRule) *InstrumentationLibraryConfigApplyConfiguration {
	b.DbStatementPayloadCollection = &value
	return b
}
