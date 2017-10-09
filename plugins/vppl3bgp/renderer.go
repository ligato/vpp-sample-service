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

//Package vppl3bgp contains a BGP VPP Agent Writer implementation for L3 VPP Plugin
package vppl3bgp

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins/localclient"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l3plugin/model/l3"
)

const description = "configuration used for Ligato VPP BGP"

// SendStaticRouteToVPP send BGP information translated to L3 default plugin structures to VPP.
func SendStaticRouteToVPP(info *bgp.ReachableIPRoute, pluginName core.PluginName) error {
	return localclient.DataChangeRequest(pluginName).Put().StaticRoute(Translate(info)).Send().ReceiveReply()
}

// Translate translates bgp information from BGP-Agent API to VPP-Agent API.
func Translate(info *bgp.ReachableIPRoute) *l3.StaticRoutes_Route {
	return &l3.StaticRoutes_Route{
		VrfId:       0,
		DstIpAddr:   info.Prefix,
		NextHopAddr: info.Nexthop.String(),
		Description: description,
		Weight:      1,
		Preference:  0}
}
