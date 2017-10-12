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

//Package vppl3bgp_test contains Ligato BGP To VPP Plugin  helper test functions
package vppl3bgp_test

import (
	"errors"
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/flavors/local"
	"github.com/ligato/cn-infra/logging/logroot"
	"github.com/ligato/vpp-sample-service/plugins/vppl3bgp"
	. "github.com/onsi/gomega"
	"net"
	"testing"
	"time"
)

const (
	prefix                = "101.0.10.0/24"
	As                    = uint32(65000)
	bgptol3testPluginName = "bgptol3testPluginName"
)

// TestHelper allows tests to be written in given/when/then idiom of BDD
type TestHelper struct {
	vars          *Variables
	Given         Given
	When          When
	Then          Then
	golangTesting *testing.T
}

// Variables is container of variables that should be accessible from every BDD component
type Variables struct {
	golangT               *testing.T
	watcher               watcher
	vppl3bgp              core.NamedPlugin
	renderer              func(*bgp.ReachableIPRoute)
	reg                   watchRegistration
	route                 *bgp.ReachableIPRoute
	lifecycleCloseChannel chan struct{}
}

// Given is composition of multiple test step methods (see BDD Given keyword)
type Given struct {
	vars *Variables
}

// When is composition of multiple test step methods (see BDD When keyword)
type When struct {
	vars *Variables
}

// Then is composition of multiple test step methods (see BDD Then keyword)
type Then struct {
	vars *Variables
}

// DefaultSetup setups needed variables and ensures that these variables are accessible from all test BDD
// components (given, when, then)
func (t *TestHelper) DefaultSetup() {
	// creating and linking variables to test parts
	t.vars = &Variables{watcher: watcher{}}
	t.Given.vars = t.vars
	t.When.vars = t.vars
	t.Then.vars = t.vars
	t.vars.golangT = t.golangTesting

	// registering gomega
	RegisterTestingT(t.vars.golangT)

	t.vars.renderer = func(route *bgp.ReachableIPRoute) {
		t.vars.route = route
	}
}

// AddNewRoute adds constant-based route to watcher callback.
func (w *When) AddNewRoute() {
	w.vars.watcher.callback(&bgp.ReachableIPRoute{
		As:      As,
		Prefix:  prefix,
		Nexthop: net.IP{101, 0, 10, 1},
	})
}

// WatcherReceivesAddedRoute checks it for route is advertized.
func (t *Then) WatcherReceivesAddedRoute() {
	Expect(t.vars.route.As).To(Equal(As))
	Expect(t.vars.route.Nexthop).To(Equal(net.IP{101, 0, 10, 1}))
	Expect(t.vars.route.Prefix).To(Equal(prefix))
}

// StopAgentLifecycle command to stop agent lifecycle.
func (w *When) StopAgentLifecycle() {
	close(w.vars.lifecycleCloseChannel)
	time.Sleep(2 * time.Second)
}

// WatcherRegistrationIsClosed assert that watcher registration Close method has been called.
func (t *Then) WatcherRegistrationIsClosed() {
	Expect(t.vars.reg.watcher_impl).To(BeNil())
}

// WatcherRegistrationIsClosed assert that watcher registration Close method has been called.
func (t *Then) Vppl3bgpPluginInitFails() {
	//starting it using cn-infra agent
	agent := core.NewAgent(logroot.StandardLogger(), 1*time.Minute,
		[]*core.NamedPlugin{&t.vars.vppl3bgp}...)
	t.vars.lifecycleCloseChannel = make(chan struct{}, 1)
	Expect(core.EventLoopWithInterrupt(agent, t.vars.lifecycleCloseChannel)).NotTo(BeNil())
}

// Vppl3bgpPluginCloseFails creates Vppl3bgpPlugin (with Watcher and Renderer registered in it) and prepares it for usage.
func (g *Given) Vppl3bgpPluginCloseFails() {

	//starting it using cn-infra agent
	g.startPluginLifecycleWithExpectedFailingClose()
}

//Vppl3bgpPluginDefault instantiate Vppl3bgp plugin with default renderer
func (g *Given) Vppl3bgpPluginDefault() {
	flavor := &local.FlavorLocal{}
	deps := *flavor.InfraDeps(bgptol3testPluginName)

	var w bgp.Watcher
	w = &g.vars.watcher
	vppl3bgpPlugin := vppl3bgp.New(vppl3bgp.Deps{
		PluginInfraDeps: deps,
		Watcher:         w,
	})
	g.vars.vppl3bgp = vppl3bgpPlugin
}

//Vppl3bgpPluginDefault instantiate Vppl3bgp plugin with injected renderer
func (g *Given) Vppl3bgpPlugin() {
	flavor := &local.FlavorLocal{}
	deps := *flavor.InfraDeps(bgptol3testPluginName)

	var w bgp.Watcher
	w = &g.vars.watcher
	vppl3bgpPlugin := vppl3bgp.New(vppl3bgp.Deps{
		PluginInfraDeps: deps,
		Watcher:         w,
		Renderer:        g.vars.renderer,
	})
	g.vars.vppl3bgp = vppl3bgpPlugin
}

// startPluginLifecycle creates cn-infra agent and uses it to start lifecycle for vppl3bgp plugin.
// expected error when close
func (g *Given) startPluginLifecycleWithExpectedFailingClose() {
	agent := core.NewAgent(logroot.StandardLogger(), 1*time.Minute,
		[]*core.NamedPlugin{&g.vars.vppl3bgp}...)
	g.vars.lifecycleCloseChannel = make(chan struct{}, 1)
	go func() {
		Expect(core.EventLoopWithInterrupt(agent, g.vars.lifecycleCloseChannel)).NotTo(BeNil())
	}()
	time.Sleep(1 * time.Millisecond)
}

// startPluginLifecycle creates cn-infra agent and uses it to start lifecycle for vppl3bgp plugin.
func (g *Given) StartPluginLifecycle() {
	agent := core.NewAgent(logroot.StandardLogger(), 1*time.Minute,
		[]*core.NamedPlugin{&g.vars.vppl3bgp}...)
	g.vars.lifecycleCloseChannel = make(chan struct{}, 1)
	go func() {
		Expect(core.EventLoopWithInterrupt(agent, g.vars.lifecycleCloseChannel)).
			To(BeNil(), "Agent's lifecycle didn't ended properly")
	}()
	time.Sleep(1 * time.Millisecond)
}

// WatcherRegistrationThrowError assigns an error to be returned by watcher when trying to register
func (g *Given) WatcherRegistrationThrowError() {
	g.vars.watcher.reg_error = errors.New("force registration error")
}

// WatcherRegistrationThrowError assigns an error to be returned by watcher when trying to close registration ticket
func (g *Given) WatcherCloseRegistrationThrowError() {
	g.vars.watcher.reg_close_error = errors.New("force registration error")
}

//watcher implementation for test purpose
type watcher struct {
	callback        func(*bgp.ReachableIPRoute)
	reg_error       error
	reg_close_error error
}

//watcher registration implementation for test purpose
type watchRegistration struct {
	watcher_impl *watcher
}

//Close ends the agreement between Plugin and watcher. Sets registration as closed
func (wr *watchRegistration) Close() error {
	err := wr.watcher_impl.reg_close_error
	wr.watcher_impl = nil
	return err
}

//WatchIPRoutes register watcher to notifications for any new learned IP-based routes.
func (watcherImpl *watcher) WatchIPRoutes(watcher string, callback func(*bgp.ReachableIPRoute)) (bgp.WatchRegistration, error) {
	watcherImpl.callback = callback
	return &watchRegistration{watcher_impl: watcherImpl}, watcherImpl.reg_error
}
