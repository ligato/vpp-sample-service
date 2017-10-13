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

// Package examples provide examples of Vpp sample service usage.
package main

import (
	"github.com/ligato/bgp-agent/bgp/gobgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins/localclient"
	local_flavor "github.com/ligato/vpp-agent/flavors/local"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/ifplugin/model/interfaces"
	"github.com/ligato/vpp-sample-service/plugins/vppl3bgp"
	"github.com/osrg/gobgp/config"
	"os"
	"time"
)

var (
	// memif1AsMaster is an a memory interface configuration used as next hop for sent BGP information.
	memif1AsMaster = interfaces.Interfaces_Interface{
		Name:    "memif1",
		Type:    interfaces.InterfaceType_MEMORY_INTERFACE,
		Enabled: true,
		Memif: &interfaces.Interfaces_Interface_Memif{
			Id:             1,
			Master:         true,
			SocketFilename: "/tmp/memif1.sock",
		},
		Mtu:         1500,
		IpAddresses: []string{"101.0.10.0/24"},
	}
	flavor = &local.FlavorLocal{}

	goBgpConfig = &config.Bgp{
		Global: config.Global{
			Config: config.GlobalConfig{
				As:       65000,
				RouterId: "172.18.0.1",
				Port:     -1,
			},
		},
		Neighbors: []config.Neighbor{
			{
				Config: config.NeighborConfig{
					PeerAs:          65001,
					NeighborAddress: "172.18.0.2",
				},
			},
		},
	}
)

const (
	bgptol3PluginName = "bgptol3example"
)

// main runs end-to-end example that demonstrates sending prefix/nexthop information from route reflector to vpp
func main() {
	deps := *flavor.InfraDeps(bgptol3PluginName)
	deps.Log.SetLevel(logging.DebugLevel)

	pluginInterface := &core.NamedPlugin{
		PluginName: bgptol3PluginName,
		Plugin:     &pluginVPPInterface{deps},
	}

	goBgpPlugin := gobgp.New(gobgp.Deps{
		PluginInfraDeps: deps,
		SessionConfig:   goBgpConfig})

	bgptol3 := vppl3bgp.New(vppl3bgp.Deps{
		PluginInfraDeps: deps,
		Watcher:         goBgpPlugin,
	})

	goBgpPluginCoreNamed := core.NamedPlugin{
		PluginName: goBgpPlugin.PluginName,
		Plugin:     goBgpPlugin}

	agent := core.NewAgent(logroot.StandardLogger(), 4*time.Minute, append(getVPPplugins(), pluginInterface,
		&bgptol3, &goBgpPluginCoreNamed)...)

	err := core.EventLoopWithInterrupt(agent, nil)
	if err != nil {
		os.Exit(1)
	}
}

//pluginVPPInterface will create required memif for VPP
type pluginVPPInterface struct {
	local.PluginInfraDeps
}

// Init creates initial structures inside VPP that are needed for prefix/next hop information sending.
// In case of problems with memif creation error informing the issue will be returned, otherwise nil will be returned
func (plugin *pluginVPPInterface) Init() error {
	return localclient.DataResyncRequest(bgptol3PluginName).Interface(&memif1AsMaster).Send().ReceiveReply()
}

// creates deaults VPP plugins
func getVPPplugins() []*core.NamedPlugin {
	flavour := local_flavor.FlavorVppLocal{}

	// no operation writer that helps avoiding NIL pointer based segmentation fault
	// used as default if some dependency was not injected
	noopWriter := &datasync.CompositeKVProtoWriter{Adapters: []datasync.KeyProtoValWriter{}}
	flavour.VPP.Publish = noopWriter
	flavour.VPP.PublishStatistics = noopWriter
	flavour.VPP.IfStatePub = noopWriter
	return flavour.Plugins()
}
