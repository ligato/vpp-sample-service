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

//Package l3writer contains a BGP VPP Agent Writer implementation for L3 VPP Plugin
package l3writer

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins/localclient"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/ifplugin/model/interfaces"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l3plugin/model/l3"
	"github.com/ligato/vpp-sample-service/bgptol3plugin/writer"
)

const description = "test"

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
		IpAddresses: []string{"192.168.1.1/24"},
	}
)

type writerImpl struct {
	pluginName core.PluginName
	writer.Writer
	dataConnector l3DataConnector
}

// New creates new L3 VPP Plugin Writer implementation
func New(pluginName core.PluginName) writer.Writer {
	return NewInjectable(pluginName, &defaultL3DataConnectorImpl{})
}

// NewInjectable creates new L3 VPP Plugin Writer implementation with a specific l3DataConnector implementation
func NewInjectable(pluginName core.PluginName, connector l3DataConnector) writer.Writer {
	return &writerImpl{
		pluginName:    pluginName,
		dataConnector: connector,
	}
}

// SendStaticRouteToVPP send BGP information translated to L3 default plugin structures to VPP.
func (writer *writerImpl) SendStaticRouteToVPP(info *bgp.ReachableIPRoute) error {
	return writer.dataConnector.DataChangeRequest(writer.pluginName).Put().StaticRoute(writer.translate(info)).Send().ReceiveReply()
}

// PrepareVPP creates initial structures inside VPP that are needed for prefix/next hop information sending.
func (writer *writerImpl) PrepareVPP() error {
	return writer.dataConnector.DataResyncRequest(writer.pluginName).Interface(&memif1AsMaster).Send().ReceiveReply()
}

// translate translates bgp information from BGP-Agent API to VPP-Agent API.
func (writer *writerImpl) translate(info *bgp.ReachableIPRoute) *l3.StaticRoutes_Route {
	return Translate(info)
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

// l3DataConnector enables data forwarding to l3 default plugin in vpp-agent
// It is interface version of github.com/ligato/vpp-agent/clientv1/defaultplugins/localclient/localclient_api.go. Main
// reason for creating this interface is ability to test functionality because mentioned localclient_api is not mockable
// (direct functions in package).
type l3DataConnector interface {
	// DataResyncRequest allows to send RESYNC requests conveniently
	DataResyncRequest(caller core.PluginName) defaultplugins.DataResyncDSL

	// DataChangeRequest allows to send Data Change requests conveniently
	DataChangeRequest(caller core.PluginName) defaultplugins.DataChangeDSL
}

// defaultL3DataConnectorImpl is implementation of l3DataConnector interface that directly connects data forwarding to default l3 default plugin
type defaultL3DataConnectorImpl struct {
	l3DataConnector
}

// DataResyncRequest allows to send RESYNC requests conveniently
func (con *defaultL3DataConnectorImpl) DataResyncRequest(caller core.PluginName) defaultplugins.DataResyncDSL {
	return localclient.DataResyncRequest(caller)
}

// DataChangeRequest allows to send Data Change requests conveniently
func (con *defaultL3DataConnectorImpl) DataChangeRequest(caller core.PluginName) defaultplugins.DataChangeDSL {
	return localclient.DataChangeRequest(caller)
}
