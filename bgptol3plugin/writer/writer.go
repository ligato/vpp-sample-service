// Copyright (c) 2017 Pantheon technologies s.r.o.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package writer contains definitions for writers to default plugin structures to VPPs (e.g. L3 Plugin)
package writer

import (
	"github.com/ligato/bgp-agent/bgp"
)

// Writer writes to default plugin structures to VPPs (e.g. L3 Plugin)
type Writer interface {
	// SendStaticRouteToVPP send BGP information translated to default plugin structures to VPP.
	SendStaticRouteToVPP(info *bgp.ReachableIPRoute) error
	// PrepareVPP creates initial structures inside VPP that are needed for prefix/next hop information sending.
	PrepareVPP() error
}
