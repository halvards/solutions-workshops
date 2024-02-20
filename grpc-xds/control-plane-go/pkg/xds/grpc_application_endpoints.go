// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xds

import (
	"slices"
	"strings"
)

type GRPCApplicationEndpoints struct {
	Node      string
	Zone      string
	Addresses []string
}

func NewGRPCApplicationEndpoints(node string, zone string, addresses []string) GRPCApplicationEndpoints {
	addressesCopy := make([]string, len(addresses))
	copy(addressesCopy, addresses)
	slices.Sort(addressesCopy)
	return GRPCApplicationEndpoints{
		Node:      node,
		Zone:      zone,
		Addresses: addressesCopy,
	}
}

// Compare assumes that the list of addresses is sorted,
// as done in `NewGRPCApplicationEndpoints()`.
func (e GRPCApplicationEndpoints) Compare(f GRPCApplicationEndpoints) int {
	if e.Node != f.Node {
		return strings.Compare(e.Node, f.Node)
	}
	if e.Zone != f.Zone {
		return strings.Compare(e.Zone, f.Zone)
	}
	return slices.Compare(e.Addresses, f.Addresses)
}

// Equal assumes that the list of addresses is sorted,
// as done in `NewGRPCApplicationEndpoints()`.
func (e GRPCApplicationEndpoints) Equal(f GRPCApplicationEndpoints) bool {
	return e.Compare(f) == 0
}
