/*
Copyright The Kubernetes Authors.

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
	v1alpha1 "github.com/vitu1234/kluster/pkg/apis/vitu.dev/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// KlusterLister helps list Klusters.
// All objects returned here must be treated as read-only.
type KlusterLister interface {
	// List lists all Klusters in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Kluster, err error)
	// Klusters returns an object that can list and get Klusters.
	Klusters(namespace string) KlusterNamespaceLister
	KlusterListerExpansion
}

// klusterLister implements the KlusterLister interface.
type klusterLister struct {
	indexer cache.Indexer
}

// NewKlusterLister returns a new KlusterLister.
func NewKlusterLister(indexer cache.Indexer) KlusterLister {
	return &klusterLister{indexer: indexer}
}

// List lists all Klusters in the indexer.
func (s *klusterLister) List(selector labels.Selector) (ret []*v1alpha1.Kluster, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Kluster))
	})
	return ret, err
}

// Klusters returns an object that can list and get Klusters.
func (s *klusterLister) Klusters(namespace string) KlusterNamespaceLister {
	return klusterNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KlusterNamespaceLister helps list and get Klusters.
// All objects returned here must be treated as read-only.
type KlusterNamespaceLister interface {
	// List lists all Klusters in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Kluster, err error)
	// Get retrieves the Kluster from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Kluster, error)
	KlusterNamespaceListerExpansion
}

// klusterNamespaceLister implements the KlusterNamespaceLister
// interface.
type klusterNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Klusters in the indexer for a given namespace.
func (s klusterNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Kluster, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Kluster))
	})
	return ret, err
}

// Get retrieves the Kluster from the indexer for a given namespace and name.
func (s klusterNamespaceLister) Get(name string) (*v1alpha1.Kluster, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("kluster"), name)
	}
	return obj.(*v1alpha1.Kluster), nil
}
