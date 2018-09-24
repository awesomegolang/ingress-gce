/*
Copyright 2015 The Kubernetes Authors.

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

package backends

import (
	computebeta "google.golang.org/api/compute/v0.beta"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/ingress-gce/pkg/composite"
	"k8s.io/ingress-gce/pkg/utils"
	"k8s.io/kubernetes/pkg/cloudprovider/providers/gce/cloud/meta"
)

// GroupKey represents a single group for a backend. The implementation of
// this Group could be InstanceGroup or NEG.
type GroupKey struct {
	Zone string
	Name string
}

// Pool is an interface to perform CRUD operations on a pool of GCE
// Backend Services.
type Pool interface {
	// Get a composite BackendService given a required version.
	Get(name string, version meta.Version) (*composite.BackendService, error)
	// Create a composite BackendService and returns it.
	Create(sp utils.ServicePort, hcLink string) (*composite.BackendService, error)
	// Update a BackendService given the composite type.
	Update(be *composite.BackendService) error
	// Delete a BackendService given its name.
	Delete(name string) error
	// Get the health of a BackendService given its name.
	Health(name string) string
	// Get a list of BackendService names currently in this pool.
	GetLocalSnapshot() []string
}

// Syncer is an interface to sync Kubernetes services to GCE BackendServices.
type Syncer interface {
	// Init an implementation of ProbeProvider.
	Init(p ProbeProvider)
	// Sync a BackendService. Implementations should only create the BackendService
	// but not its groups.
	Sync(svcPorts []utils.ServicePort) error
	// GC garbage collects unused BackendService's
	GC(svcPorts []utils.ServicePort) error
	// Status returns the status of a BackendService given its name.
	Status(name string) string
	// Shutdown cleans up all BackendService's previously synced.
	Shutdown() error
}

// Linker is an interface to link backends with their associated groups.
type Linker interface {
	// Link a BackendService to its groups.
	Link(sp utils.ServicePort, groups []GroupKey) error
}

// NEGGetter is an interface to retrieve NEG object
type NEGGetter interface {
	GetNetworkEndpointGroup(name string, zone string) (*computebeta.NetworkEndpointGroup, error)
}

// ProbeProvider retrieves a probe struct given a nodePort
type ProbeProvider interface {
	GetProbe(sp utils.ServicePort) (*api_v1.Probe, error)
}
