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

//Package vppl3bgp_test contains Ligato BGP To VPP Plugin implementation tests
package vppl3bgp_test

import (
	"testing"
)

func TestVppl3bgpPlugin(x *testing.T) {
	t := TestHelper{golangTesting: x}
	t.DefaultSetup()
	t.Given.Vppl3bgpPlugin()
	t.Given.StartPluginLifecycle()
	t.When.AddNewRoute()
	t.Then.WatcherReceivesAddedRoute()
	t.When.StopAgentLifecycle()
	t.Then.WatcherRegistrationIsClosed()
}

func TestVppl3bgpPluginDefault(x *testing.T) {
	t := TestHelper{golangTesting: x}
	t.DefaultSetup()
	t.Given.Vppl3bgpPluginDefault()
	t.Given.StartPluginLifecycle()
	t.When.AddNewRoute()
	t.When.StopAgentLifecycle()
	t.Then.WatcherRegistrationIsClosed()
}

func TestVppl3bgpPluginInitFail(x *testing.T) {
	t := TestHelper{golangTesting: x}
	t.DefaultSetup()
	t.Given.WatcherRegistrationThrowError()
	t.Given.Vppl3bgpPlugin()
	t.Then.Vppl3bgpPluginInitFails()
}

func TestVppl3bgpPluginCloseFail(x *testing.T) {
	t := TestHelper{golangTesting: x}
	t.DefaultSetup()
	t.Given.WatcherCloseRegistrationThrowError()
	t.Given.Vppl3bgpPlugin()
	t.Given.startPluginLifecycleWithExpectedFailingClose()
	t.When.AddNewRoute()
	t.Then.WatcherReceivesAddedRoute()
	t.When.StopAgentLifecycle()
}
