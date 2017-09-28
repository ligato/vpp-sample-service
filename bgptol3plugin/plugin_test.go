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

//Package agent_test contains BGPtoL3 plugin implementation tests
package bgptol3plugin_test

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/vpp-sample-service/bgptol3plugin"
	assertProvider "github.com/stretchr/testify/assert"
	"net"
	errors "pantheon.tech/ligato-bgp/agent/utils/error"
	"pantheon.tech/ligato-bgp/bgp-vpp-agent/bgptol3plugin/writer"
	"testing"
	"time"
)

var (
	route = bgp.ReachableIPRoute{As: 1, Prefix: "1.2.3.4/32", Nexthop: net.IPv4(192, 168, 1, 1)}
)

// TestPlugin test BGPtoL3 plugin for correct receiving and forwarding BGP information
func TestPlugin(t *testing.T) {
	// prepare of tested instances/helper instances
	assert := assertProvider.New(t)
	bgpChannel := make(chan bgp.ReachableIPRoute, 1)
	mockWriter := mockWriter{}
	bgpVppAgent := bgptol3plugin.NewInjectable(bgpChannel, &mockWriter)

	//initialize plugin
	errors.PanicIfError(bgpVppAgent.Init())
	if postInitPlugin, ok := bgpVppAgent.Plugin.(core.PostInit); ok {
		errors.PanicIfError(postInitPlugin.AfterInit())
	}

	//sent something to plugin from BGP-Agent side (and give plugin time to process, possibly in other go routines)
	bgpChannel <- route
	time.Sleep(500 * time.Millisecond) //this should be long enough for plugin to forward the information to l3 default plugin

	//checking catched objects sent to mockWriter (not to l3 default plugin directly)
	assert.NotNil(mockWriter.catchedStaticRoute)
	assert.Equal(route.Prefix, mockWriter.catchedStaticRoute.Prefix)
	assert.Equal(route.Nexthop, mockWriter.catchedStaticRoute.Nexthop)
	assert.True(mockWriter.vppPrepared)

	//close plugin
	errors.PanicIfError(bgpVppAgent.Close())
}

type mockWriter struct {
	writer.Writer
	catchedStaticRoute *bgp.ReachableIPRoute
	vppPrepared        bool
}

// SendStaticRouteToVPP send BGP information translated to L3 default plugin structures to VPP.
func (writer *mockWriter) SendStaticRouteToVPP(info *bgp.ReachableIPRoute) error {
	writer.catchedStaticRoute = info
	return nil
}

// PrepareVPP creates initial structures inside VPP that are needed for prefix/next hop information sending.
func (writer *mockWriter) PrepareVPP() error {
	writer.vppPrepared = true
	return nil
}
