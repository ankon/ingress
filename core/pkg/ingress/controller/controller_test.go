/*
Copyright 2017 The Kubernetes Authors.

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

package controller

import (
	"reflect"
	"testing"

	"k8s.io/ingress/core/pkg/ingress"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/util/sets"
)

func TestGetDefaultUpstream(t *testing.T) {
	const ds = "default-service"
	svcStore := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{})
	ids := sets.NewString(ds)
	for id := range ids {
		svcStore.Add(&api.Service{ObjectMeta: api.ObjectMeta{Name: id}})
	}

	epStore := cache.NewStore(cache.MetaNamespaceKeyFunc)
	ic := &GenericController{}
	ic.svcLister = cache.StoreToServiceLister{Indexer: svcStore}
	ic.endpLister = cache.StoreToEndpointsLister{Store: epStore}
	ic.cfg = &Configuration{
		DefaultService: ds,
	}
	db := ic.getDefaultUpstream()
	if !reflect.DeepEqual(db.Endpoints, []ingress.Endpoint{newDefaultServer()}) {
		t.Errorf("expected default backend (%v) but %v was returned", db.Endpoints, db)
	}
}
