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

//Package l3writer_test contains l3writer implementation tests
package l3writer_test

import (
	"github.com/golang/mock/gomock"
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/vpp-agent/clientv1/defaultplugins"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/ifplugin/model/interfaces"
	"github.com/ligato/vpp-agent/plugins/defaultplugins/l3plugin/model/l3"
	"github.com/ligato/vpp-sample-service/bgptol3plugin/writer/l3writer"
	"github.com/ligato/vpp-sample-service/mocks"
	assertProvider "github.com/stretchr/testify/assert"
	"net"
	"testing"
)

const testPluginName core.PluginName = "testPlugin"

var (
	route = bgp.ReachableIPRoute{As: 1, Prefix: "1.2.3.4/32", Nexthop: net.IPv4(192, 168, 1, 1)}
)

// TestL3Writer test l3Writer implementation if it uses properly l3 default plugin API
func TestL3Writer(t *testing.T) {
	// prepare of tested instances/helper instances
	assert := assertProvider.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	conMock := mockL3DataConnector{mockCtrl: mockCtrl}
	writer := l3writer.NewInjectable(testPluginName, &conMock)

	// check VPP prepare
	writer.PrepareVPP()
	assert.NotNil(conMock.CatchedVppInterface)
	assert.True(conMock.CatchedVppInterfaceSent)

	// check static route write
	writer.SendStaticRouteToVPP(&route)
	assert.NotNil(conMock.CatchedStaticRoute)
	assert.Equal(route.Prefix, conMock.CatchedStaticRoute.DstIpAddr)
	assert.Equal(route.Nexthop.String(), conMock.CatchedStaticRoute.NextHopAddr)
	assert.True(conMock.CatchedStaticRouteSent)
}

type mockL3DataConnector struct {
	mockCtrl                *gomock.Controller
	CatchedStaticRoute      *l3.StaticRoutes_Route
	CatchedStaticRouteSent  bool
	CatchedVppInterface     *interfaces.Interfaces_Interface
	CatchedVppInterfaceSent bool
}

// DataResyncRequest allows to send RESYNC requests conveniently
func (con *mockL3DataConnector) DataResyncRequest(caller core.PluginName) defaultplugins.DataResyncDSL {
	resyncDSLMock := mocks.NewMockDataResyncDSL(con.mockCtrl)
	replyMock := mocks.NewMockReply(con.mockCtrl)

	gomock.InOrder(
		resyncDSLMock.EXPECT().Interface(gomock.Any()).Do(func(intf *interfaces.Interfaces_Interface) {
			con.CatchedVppInterface = intf
		}).Return(resyncDSLMock),
		resyncDSLMock.EXPECT().Send().Do(func() {
			con.CatchedVppInterfaceSent = true
		}).Return(replyMock),
		replyMock.EXPECT().ReceiveReply().Return(nil),
	)

	return resyncDSLMock
}

// DataChangeRequest allows to send Data Change requests conveniently
func (con *mockL3DataConnector) DataChangeRequest(caller core.PluginName) defaultplugins.DataChangeDSL {
	changeDSLMock := mocks.NewMockDataChangeDSL(con.mockCtrl)
	putDSLMock := mocks.NewMockPutDSL(con.mockCtrl)
	replyMock := mocks.NewMockReply(con.mockCtrl)

	gomock.InOrder(
		changeDSLMock.EXPECT().Put().Return(putDSLMock),
		putDSLMock.EXPECT().StaticRoute(gomock.Any()).Do(func(route *l3.StaticRoutes_Route) {
			con.CatchedStaticRoute = route
		}).Return(putDSLMock),
		putDSLMock.EXPECT().Send().Do(func() {
			con.CatchedStaticRouteSent = true
		}).Return(replyMock),
		replyMock.EXPECT().ReceiveReply().Return(nil),
	)

	return changeDSLMock
}
