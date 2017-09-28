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

// Package examples provide examples of BGP-VPP-Agent usage.
package main

import (
	"fmt"
	"github.com/ligato/bgp-agent/bgp"
	ligatoAgent "github.com/ligato/cn-infra/core"
	"github.com/ligato/cn-infra/logging"
	log "github.com/ligato/cn-infra/logging/logrus"
	local_flavor "github.com/ligato/vpp-agent/flavors/local"
	"github.com/ligato/vpp-sample-service/bgptol3plugin"
	"os"
	"pantheon.tech/ligato-bgp/agent/bgp/gobgp"
	gobgpconfigs "pantheon.tech/ligato-bgp/agent/bgp/gobgp/configs"
	"time"
)

const (
	// exampleDuration provides duration of example and by using it we have the ability to show flawless and graceful stop of all components
	exampleDuration = 4 * time.Minute

	// minimalExampleDuration provides recommended duration of this example to get expected behaviour.
	// GoBGP can take around a minute to properly initialize so we take 2 minutes to be sure. After that information can pass
	// and it should take probably no more than 10 seconds.
	minimalExampleDuration = 2*time.Minute + 10*time.Second
)

// init sets the default logging level
func init() {
	log.DefaultLogger().SetOutput(os.Stdout)
	log.DefaultLogger().SetLevel(logging.DebugLevel)
}

// main runs end-to-end example that demonstrates sending prefix/nexthop information from route reflector to vpp
func main() {
	// creating connection channel between bgp agent and bgp-vpp agent (and channel to stop them)
	connectionChannel := make(chan bgp.ReachableIPRoute, 10)
	bgpAgentStopChannel := make(chan struct{})
	vppAgentStopChannel := make(chan struct{})
	exampleStopChannel := make(chan struct{})
	logger := log.DefaultLogger()

	// running agents in separate go routines
	go runVPPAgentWithBGPtoL3Plugin(connectionChannel, vppAgentStopChannel, exampleStopChannel, logger)
	go runBgpAgent(connectionChannel, bgpAgentStopChannel, vppAgentStopChannel, exampleStopChannel, logger)

	//stopping example run
	stopAllAgentsIn(exampleDuration, bgpAgentStopChannel, exampleStopChannel, logger)
	logger.Info("Example finished.")
}

// runBgpAgent starts the BGP-Agent.
func runBgpAgent(connectionChannel chan bgp.ReachableIPRoute, bgpStopChannel chan struct{}, vppStopChannel chan struct{},
	exampleStopChannel chan struct{}, logger logging.Logger) {
	// creation of BGP agent
	goBgpPlugin := gobgp.New(gobgpconfigs.AllInDocker)
	bgpAgentVar, err := bgpAgent.New([]*bgp.Plugin{&goBgpPlugin.Plugin})
	if err != nil {
		logger.Panic("BGP-Agent can't be created: %v", err)
	}
	logger.Debug("BGP-Agent successfully created.")

	// watch for new BGP data
	if err := bgpAgentVar.WatchIPRoutes("BGP-VPP Ligato plugin", bgp.ToChan(connectionChannel, logger)); err != nil {
		logger.Panicf("Unable to watch BGP agent for learned BGP data: %v", err)
	}

	// start BGP-Agent and its plugins
	if err := startBGPAgent(logger, bgpAgentVar); err != nil {
		logger.Panicf("Unable to start BGP agent: %v", err)
	}

	// wait for end of example
	<-bgpStopChannel

	// start BGP-Agent and its plugins
	if err := stopBGPAgent(logger, bgpAgentVar); err != nil {
		logger.Panicf("Unable to stop BGP agent: %v", err)
	}

	// VPP-Agent with BGP-to-L3 plugin can be stopped only after BGP-Agent flushes all information to it
	vppStopChannel <- struct{}{}

	// signal that BGP Agent is stopped (example end don't need to wait for BGP Agent to shutdown properly anymore)
	exampleStopChannel <- struct{}{}
}

// stopBGPAgent closes BGP-Agent service and stops BGP-Agent
func stopBGPAgent(logger logging.Logger, bgpAgent bgp.Agent) (err error) {
	logger.Debug("BGP-Agent stopping...")
	if err = bgpAgent.Close(); err != nil {
		err = fmt.Errorf("BGP-Agent failed to close: %v", err)
		return
	}
	logger.Info("BGP-Agent stopped.")
	return
}

// startBGPAgent starts BGP-Agent and instantiate its service
func startBGPAgent(logger logging.Logger, bgpAgent bgp.Agent) (err error) {
	logger.Debug("Starting lifecycle of BGP-Agent...")
	if err = bgpAgent.Init(); err != nil {
		err = fmt.Errorf("BGP-Agent failed to initialize: %v", err)
		return
	}
	logger.Debug("BGP-Agent initialized.")

	return
}

// runVPPAgentWithBGPtoL3Plugin starts the VPP-Agent with the BGP-to-L3 plugin
func runVPPAgentWithBGPtoL3Plugin(bgpChannel chan bgp.ReachableIPRoute, vppAgentStopChannel chan struct{}, exampleStopChannel chan struct{}, logger logging.Logger) {
	// Create BGP-to-L3 plugin that is plugin of the Vpp Agent
	bgptol3 := bgptol3plugin.New(bgpChannel)

	// plugins set(=flavor) for local linux environment with vpp
	flavour := local_flavor.FlavorVppLocal{}

	// Create new ligato agent
	ligatoAgentVar := ligatoAgent.NewAgent(logger, 15*time.Second, append(flavour.Plugins(), &bgptol3)...)

	// Run agent in event loop
	ligatoAgent.EventLoopWithInterrupt(ligatoAgentVar, vppAgentStopChannel)

	// signal that BgpVpp Agent and Ligato is stopped (example end don't need to wait for proper shutdown anymore)
	exampleStopChannel <- struct{}{}
}

// stopAllAgentsIn stops all agents(Vpp-Agent,BGP-Agent,BGP-VPP-Agent) after given time period.
func stopAllAgentsIn(sleepTime time.Duration, bgpAgentStopChannel chan struct{}, exampleStopChannel chan struct{}, logger logging.Logger) {
	waitUntilExampleEnd(sleepTime, logger)
	stopAllAgents(bgpAgentStopChannel, exampleStopChannel, logger)
}

// stopAllAgents stops all agents(Vpp-Agent,BGP-Agent) by using channel messaging.
func stopAllAgents(bgpAgentStopChannel chan struct{}, exampleStopChannel chan struct{}, logger logging.Logger) {
	logger.Info("Starting to stop all agents.")

	// trigger stopping of BGP-Agent (in its endphase it will also trigger stopping of VPP-Agent with BGP-to-L3 plugin)
	bgpAgentStopChannel <- struct{}{}

	// wait for both agents to stop gracefully
	<-exampleStopChannel
	<-exampleStopChannel

	logger.Info("All Agents stopped.")
}

// waitUntilExampleEnd makes current goroutine sleep until given time period, but it will not awake sooner than time period given in minimalExampleDuration.
func waitUntilExampleEnd(sleepTime time.Duration, logger logging.Logger) {
	if sleepTime < minimalExampleDuration {
		sleepTime = minimalExampleDuration
		logger.Warnf("Time to stop all agents is too low to observe expected behaviour. Time is raised to %v", sleepTime)
	}
	time.Sleep(sleepTime)
}
