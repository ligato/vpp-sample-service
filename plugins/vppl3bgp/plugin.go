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

// Package vppl3bgp implements the Vpp sample service plugin that allows plugin
// to render learned IP-based routes to l3 plugin.
package vppl3bgp

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins/localclient"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l3plugin/model/l3"
)

// pluginID of BGP-to-L3 plugin
const pluginID core.PluginName = "bgp-to-l3-plugin"

// description for l3 StaticRoutes_Route
const description = "configuration used for Ligato VPP BGP"

// Plugin with BGP functionality (VPP Agent plugin that servers as BGP-VPP Agent)
// it handles information coming for BGP-Agent channel and sends them transformed to L3 default plugin.
type pluginImpl struct {
	Deps
	reg bgp.WatchRegistration
}

// Deps combines all needed dependencies for Plugin struct. These dependencies should be injected into Plugin by using constructor's Deps parameter.
type Deps struct {
	local.PluginInfraDeps //inject
	Watcher               Watcher
	Renderer              func(*bgp.ReachableIPRoute) //inject optional (mainly for testing purposes)
}

// New creates Plugin with learned IP-based route to l3 plugin rendering functionality by default.
// Renderer can be injected via Dependencies <deps>
func New(deps Deps) core.NamedPlugin {
	return core.NamedPlugin{
		PluginName: pluginID,
		Plugin: &pluginImpl{
			Deps: deps,
		},
	}
}

// Init logs attempt of plugin initialization to be sure that plugin is properly recognized. No initialization
// of plugin is not needed yet.
func (plugin *pluginImpl) Init() error {
	if plugin.Deps.Renderer == nil {
		plugin.Deps.Renderer = func(route *bgp.ReachableIPRoute) {
			plugin.Log.Debugf("SendStaticRouteToVPP %v", route)
			err := localclient.DataChangeRequest(pluginID).Put().StaticRoute(translate(route)).Send().ReceiveReply()
			if err != nil {
				plugin.Log.Errorf("Failed to send route %v to VPP. %v", route, err)
			}
		}
	}

	reg, err := plugin.Watcher.WatchIPRoutes("BGP-VPP Ligato plugin", plugin.Deps.Renderer)
	plugin.reg = reg
	plugin.Log.Info("Initialization of the BGP plugin has completed")
	return err
}

// translate translates bgp information from BGP-Agent API to VPP-Agent API.
func translate(info *bgp.ReachableIPRoute) *l3.StaticRoutes_Route {
	return &l3.StaticRoutes_Route{
		VrfId:       0,
		DstIpAddr:   info.Prefix,
		NextHopAddr: info.Nexthop.String(),
		Description: description,
		Weight:      1,
		Preference:  0}
}

//Close ends the agreement between Plugin and watcher. Plugin stops sending watcher any further notifications.
func (plugin *pluginImpl) Close() error {
	return plugin.reg.Close()
}

//Watcher common interface between Ligato BGP Plugin for register watcher to notifications
type Watcher interface {
	//WatchIPRoutes register watcher to notifications for any new learned IP-based routes.
	WatchIPRoutes(watcher string, callback func(*bgp.ReachableIPRoute)) (bgp.WatchRegistration, error)
}
