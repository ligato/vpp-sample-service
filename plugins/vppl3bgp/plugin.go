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

// Package agent provides core agent functionality (e.g.VPP-Agent plugin implementation)
package vppl3bgp

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/flavors/local"
)

// PluginID of BGP-to-L3 plugin
const PluginID core.PluginName = "bgp-to-l3-plugin"

// Plugin with BGP functionality (VPP Agent plugin that servers as BGP-VPP Agent)
// it handles information coming for BGP-Agent channel and sends them transformed to L3 default plugin.
type pluginImpl struct {
	Deps
}

type Deps struct {
	local.PluginInfraDeps //inject
	Watcher               Watcher
	Renderer              func(*bgp.ReachableIPRoute) //inject optional (mainly for testing purposes)
}

// New creates Plugin with BGP functionality
func New(deps Deps) core.NamedPlugin {
	return NewInjectable(deps)
}

// NewInjectable creates Plugin with BGP functionality with an specific writer implementation
func NewInjectable(deps Deps) core.NamedPlugin {
	pluginVar := &pluginImpl{
		Deps: deps,
	}
	return core.NamedPlugin{
		PluginName: PluginID,
		Plugin:     pluginVar,
	}
}

// Init logs attempt of plugin initialization to be sure that plugin is properly recognized. No initialization of plugin is not needed yet.
func (plugin *pluginImpl) Init() error {
	if plugin.Deps.Renderer == nil {
		plugin.Deps.Renderer = func(route *bgp.ReachableIPRoute) {
			err := SendStaticRouteToVPP(route, PluginID)
			if err != nil {
				plugin.Log.Errorf("Failed to send route %v to VPP. %v", route, err)
			}
		}
	}

	plugin.Log.Info("Initialization of the BGP plugin has completed")
	plugin.Watcher.Watch("BGP-VPP Ligato plugin", plugin.Deps.Renderer)
	return nil
}

type Watcher interface {
	Watch(subscriber string, callback func(*bgp.ReachableIPRoute))
}
