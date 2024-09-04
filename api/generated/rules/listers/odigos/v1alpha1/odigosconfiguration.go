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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// OdigosConfigurationLister helps list OdigosConfigurations.
// All objects returned here must be treated as read-only.
type OdigosConfigurationLister interface {
	// List lists all OdigosConfigurations in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.OdigosConfiguration, err error)
	// OdigosConfigurations returns an object that can list and get OdigosConfigurations.
	OdigosConfigurations(namespace string) OdigosConfigurationNamespaceLister
	OdigosConfigurationListerExpansion
}

// odigosConfigurationLister implements the OdigosConfigurationLister interface.
type odigosConfigurationLister struct {
	listers.ResourceIndexer[*v1alpha1.OdigosConfiguration]
}

// NewOdigosConfigurationLister returns a new OdigosConfigurationLister.
func NewOdigosConfigurationLister(indexer cache.Indexer) OdigosConfigurationLister {
	return &odigosConfigurationLister{listers.New[*v1alpha1.OdigosConfiguration](indexer, v1alpha1.Resource("odigosconfiguration"))}
}

// OdigosConfigurations returns an object that can list and get OdigosConfigurations.
func (s *odigosConfigurationLister) OdigosConfigurations(namespace string) OdigosConfigurationNamespaceLister {
	return odigosConfigurationNamespaceLister{listers.NewNamespaced[*v1alpha1.OdigosConfiguration](s.ResourceIndexer, namespace)}
}

// OdigosConfigurationNamespaceLister helps list and get OdigosConfigurations.
// All objects returned here must be treated as read-only.
type OdigosConfigurationNamespaceLister interface {
	// List lists all OdigosConfigurations in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.OdigosConfiguration, err error)
	// Get retrieves the OdigosConfiguration from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.OdigosConfiguration, error)
	OdigosConfigurationNamespaceListerExpansion
}

// odigosConfigurationNamespaceLister implements the OdigosConfigurationNamespaceLister
// interface.
type odigosConfigurationNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.OdigosConfiguration]
}