//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/keyval-dev/odigos/common"
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorsGroup) DeepCopyInto(out *CollectorsGroup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorsGroup.
func (in *CollectorsGroup) DeepCopy() *CollectorsGroup {
	if in == nil {
		return nil
	}
	out := new(CollectorsGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CollectorsGroup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorsGroupList) DeepCopyInto(out *CollectorsGroupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CollectorsGroup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorsGroupList.
func (in *CollectorsGroupList) DeepCopy() *CollectorsGroupList {
	if in == nil {
		return nil
	}
	out := new(CollectorsGroupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CollectorsGroupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorsGroupSpec) DeepCopyInto(out *CollectorsGroupSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorsGroupSpec.
func (in *CollectorsGroupSpec) DeepCopy() *CollectorsGroupSpec {
	if in == nil {
		return nil
	}
	out := new(CollectorsGroupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CollectorsGroupStatus) DeepCopyInto(out *CollectorsGroupStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CollectorsGroupStatus.
func (in *CollectorsGroupStatus) DeepCopy() *CollectorsGroupStatus {
	if in == nil {
		return nil
	}
	out := new(CollectorsGroupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Destination) DeepCopyInto(out *Destination) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Destination.
func (in *Destination) DeepCopy() *Destination {
	if in == nil {
		return nil
	}
	out := new(Destination)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Destination) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DestinationList) DeepCopyInto(out *DestinationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Destination, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DestinationList.
func (in *DestinationList) DeepCopy() *DestinationList {
	if in == nil {
		return nil
	}
	out := new(DestinationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DestinationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DestinationSpec) DeepCopyInto(out *DestinationSpec) {
	*out = *in
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	if in.Signals != nil {
		in, out := &in.Signals, &out.Signals
		*out = make([]common.ObservabilitySignal, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DestinationSpec.
func (in *DestinationSpec) DeepCopy() *DestinationSpec {
	if in == nil {
		return nil
	}
	out := new(DestinationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DestinationStatus) DeepCopyInto(out *DestinationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DestinationStatus.
func (in *DestinationStatus) DeepCopy() *DestinationStatus {
	if in == nil {
		return nil
	}
	out := new(DestinationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationConfig) DeepCopyInto(out *InstrumentationConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationConfig.
func (in *InstrumentationConfig) DeepCopy() *InstrumentationConfig {
	if in == nil {
		return nil
	}
	out := new(InstrumentationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstrumentationConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationConfigList) DeepCopyInto(out *InstrumentationConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstrumentationConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationConfigList.
func (in *InstrumentationConfigList) DeepCopy() *InstrumentationConfigList {
	if in == nil {
		return nil
	}
	out := new(InstrumentationConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstrumentationConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationConfigSpec) DeepCopyInto(out *InstrumentationConfigSpec) {
	*out = *in
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = make([]WorkloadInstrumentationConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationConfigSpec.
func (in *InstrumentationConfigSpec) DeepCopy() *InstrumentationConfigSpec {
	if in == nil {
		return nil
	}
	out := new(InstrumentationConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentationLibrary) DeepCopyInto(out *InstrumentationLibrary) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentationLibrary.
func (in *InstrumentationLibrary) DeepCopy() *InstrumentationLibrary {
	if in == nil {
		return nil
	}
	out := new(InstrumentationLibrary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentedApplication) DeepCopyInto(out *InstrumentedApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentedApplication.
func (in *InstrumentedApplication) DeepCopy() *InstrumentedApplication {
	if in == nil {
		return nil
	}
	out := new(InstrumentedApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstrumentedApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentedApplicationList) DeepCopyInto(out *InstrumentedApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstrumentedApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentedApplicationList.
func (in *InstrumentedApplicationList) DeepCopy() *InstrumentedApplicationList {
	if in == nil {
		return nil
	}
	out := new(InstrumentedApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstrumentedApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentedApplicationSpec) DeepCopyInto(out *InstrumentedApplicationSpec) {
	*out = *in
	if in.Languages != nil {
		in, out := &in.Languages, &out.Languages
		*out = make([]common.LanguageByContainer, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentedApplicationSpec.
func (in *InstrumentedApplicationSpec) DeepCopy() *InstrumentedApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(InstrumentedApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstrumentedApplicationStatus) DeepCopyInto(out *InstrumentedApplicationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstrumentedApplicationStatus.
func (in *InstrumentedApplicationStatus) DeepCopy() *InstrumentedApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(InstrumentedApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OdigosConfiguration) DeepCopyInto(out *OdigosConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OdigosConfiguration.
func (in *OdigosConfiguration) DeepCopy() *OdigosConfiguration {
	if in == nil {
		return nil
	}
	out := new(OdigosConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OdigosConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OdigosConfigurationList) DeepCopyInto(out *OdigosConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OdigosConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OdigosConfigurationList.
func (in *OdigosConfigurationList) DeepCopy() *OdigosConfigurationList {
	if in == nil {
		return nil
	}
	out := new(OdigosConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OdigosConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OdigosConfigurationSpec) DeepCopyInto(out *OdigosConfigurationSpec) {
	*out = *in
	if in.IgnoredNamespaces != nil {
		in, out := &in.IgnoredNamespaces, &out.IgnoredNamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.SupportedSDKs != nil {
		in, out := &in.SupportedSDKs, &out.SupportedSDKs
		*out = make(map[common.ProgrammingLanguage][]SupportedOtelSdk, len(*in))
		for key, val := range *in {
			var outVal []SupportedOtelSdk
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]SupportedOtelSdk, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.DefaultSDKs != nil {
		in, out := &in.DefaultSDKs, &out.DefaultSDKs
		*out = make(map[common.ProgrammingLanguage]SupportedOtelSdk, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OdigosConfigurationSpec.
func (in *OdigosConfigurationSpec) DeepCopy() *OdigosConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(OdigosConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SupportedOtelSdk) DeepCopyInto(out *SupportedOtelSdk) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SupportedOtelSdk.
func (in *SupportedOtelSdk) DeepCopy() *SupportedOtelSdk {
	if in == nil {
		return nil
	}
	out := new(SupportedOtelSdk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadInstrumentationConfig) DeepCopyInto(out *WorkloadInstrumentationConfig) {
	*out = *in
	if in.InstrumentationLibraries != nil {
		in, out := &in.InstrumentationLibraries, &out.InstrumentationLibraries
		*out = make([]InstrumentationLibrary, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadInstrumentationConfig.
func (in *WorkloadInstrumentationConfig) DeepCopy() *WorkloadInstrumentationConfig {
	if in == nil {
		return nil
	}
	out := new(WorkloadInstrumentationConfig)
	in.DeepCopyInto(out)
	return out
}
