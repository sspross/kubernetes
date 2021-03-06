/*
Copyright 2015 Google Inc. All rights reserved.

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

package client

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

type SecretsNamespacer interface {
	Secrets(namespace string) SecretsInterface
}

type SecretsInterface interface {
	Create(secret *api.Secret) (*api.Secret, error)
	Update(secret *api.Secret) (*api.Secret, error)
	Delete(name string) error
	List(label labels.Selector, field fields.Selector) (*api.SecretList, error)
	Get(name string) (*api.Secret, error)
	Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
}

// events implements Secrets interface
type secrets struct {
	client    *Client
	namespace string
}

// newSecrets returns a new secrets object.
func newSecrets(c *Client, ns string) *secrets {
	return &secrets{
		client:    c,
		namespace: ns,
	}
}

func (s *secrets) Create(secret *api.Secret) (*api.Secret, error) {
	result := &api.Secret{}
	err := s.client.Post().
		Namespace(secret.Namespace).
		Resource("secrets").
		Body(secret).
		Do().
		Into(result)

	return result, err
}

// List returns a list of secrets matching the selectors.
func (s *secrets) List(label labels.Selector, field fields.Selector) (*api.SecretList, error) {
	result := &api.SecretList{}

	err := s.client.Get().
		Namespace(s.namespace).
		Resource("secrets").
		LabelsSelectorParam(api.LabelSelectorQueryParam(s.client.APIVersion()), label).
		FieldsSelectorParam(api.FieldSelectorQueryParam(s.client.APIVersion()), field).
		Do().
		Into(result)

	return result, err
}

// Get returns the given secret, or an error.
func (s *secrets) Get(name string) (*api.Secret, error) {
	result := &api.Secret{}
	err := s.client.Get().
		Namespace(s.namespace).
		Resource("secrets").
		Name(name).
		Do().
		Into(result)

	return result, err
}

// Watch starts watching for secrets matching the given selectors.
func (s *secrets) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	return s.client.Get().
		Prefix("watch").
		Namespace(s.namespace).
		Resource("secrets").
		Param("resourceVersion", resourceVersion).
		LabelsSelectorParam(api.LabelSelectorQueryParam(s.client.APIVersion()), label).
		FieldsSelectorParam(api.FieldSelectorQueryParam(s.client.APIVersion()), field).
		Watch()
}

func (s *secrets) Delete(name string) error {
	return s.client.Delete().
		Namespace(s.namespace).
		Resource("secrets").
		Name(name).
		Do().
		Error()
}

func (s *secrets) Update(secret *api.Secret) (result *api.Secret, err error) {
	result = &api.Secret{}
	err = s.client.Put().
		Namespace(s.namespace).
		Resource("secrets").
		Name(secret.Name).
		Body(secret).
		Do().
		Into(result)

	return
}
