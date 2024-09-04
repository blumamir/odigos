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
	common "github.com/odigos-io/odigos/common"
)

// SdkConfigApplyConfiguration represents a declarative configuration of the SdkConfig type for use
// with apply.
type SdkConfigApplyConfiguration struct {
	Language                             *common.ProgrammingLanguage                      `json:"language,omitempty"`
	InstrumentationLibraryConfigs        []InstrumentationLibraryConfigApplyConfiguration `json:"instrumentationLibraryConfigs,omitempty"`
	HeadSamplingConfig                   *HeadSamplingConfigApplyConfiguration            `json:"headSamplerConfig,omitempty"`
	DefaultHttpRequestPayloadCollection  *rulesv1alpha1.HttpPayloadCollectionRule         `json:"defaultHttpRequestPayloadCollection,omitempty"`
	DefaultHttpResponsePayloadCollection *rulesv1alpha1.HttpPayloadCollectionRule         `json:"defaultHttpResponsePayloadCollection,omitempty"`
	DefaultDbQueryPayloadCollection      *rulesv1alpha1.DbQueryPayloadCollectionRule      `json:"defaultDbQueryPayloadCollection,omitempty"`
}

// SdkConfigApplyConfiguration constructs a declarative configuration of the SdkConfig type for use with
// apply.
func SdkConfig() *SdkConfigApplyConfiguration {
	return &SdkConfigApplyConfiguration{}
}

// WithLanguage sets the Language field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Language field is set to the value of the last call.
func (b *SdkConfigApplyConfiguration) WithLanguage(value common.ProgrammingLanguage) *SdkConfigApplyConfiguration {
	b.Language = &value
	return b
}

// WithInstrumentationLibraryConfigs adds the given value to the InstrumentationLibraryConfigs field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the InstrumentationLibraryConfigs field.
func (b *SdkConfigApplyConfiguration) WithInstrumentationLibraryConfigs(values ...*InstrumentationLibraryConfigApplyConfiguration) *SdkConfigApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithInstrumentationLibraryConfigs")
		}
		b.InstrumentationLibraryConfigs = append(b.InstrumentationLibraryConfigs, *values[i])
	}
	return b
}

// WithHeadSamplingConfig sets the HeadSamplingConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the HeadSamplingConfig field is set to the value of the last call.
func (b *SdkConfigApplyConfiguration) WithHeadSamplingConfig(value *HeadSamplingConfigApplyConfiguration) *SdkConfigApplyConfiguration {
	b.HeadSamplingConfig = value
	return b
}

// WithDefaultHttpRequestPayloadCollection sets the DefaultHttpRequestPayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DefaultHttpRequestPayloadCollection field is set to the value of the last call.
func (b *SdkConfigApplyConfiguration) WithDefaultHttpRequestPayloadCollection(value rulesv1alpha1.HttpPayloadCollectionRule) *SdkConfigApplyConfiguration {
	b.DefaultHttpRequestPayloadCollection = &value
	return b
}

// WithDefaultHttpResponsePayloadCollection sets the DefaultHttpResponsePayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DefaultHttpResponsePayloadCollection field is set to the value of the last call.
func (b *SdkConfigApplyConfiguration) WithDefaultHttpResponsePayloadCollection(value rulesv1alpha1.HttpPayloadCollectionRule) *SdkConfigApplyConfiguration {
	b.DefaultHttpResponsePayloadCollection = &value
	return b
}

// WithDefaultDbQueryPayloadCollection sets the DefaultDbQueryPayloadCollection field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DefaultDbQueryPayloadCollection field is set to the value of the last call.
func (b *SdkConfigApplyConfiguration) WithDefaultDbQueryPayloadCollection(value rulesv1alpha1.DbQueryPayloadCollectionRule) *SdkConfigApplyConfiguration {
	b.DefaultDbQueryPayloadCollection = &value
	return b
}