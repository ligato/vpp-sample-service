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
package bgptol3plugin

import (
	"github.com/ligato/bgp-agent/bgp"
	"github.com/ligato/cn-infra/core"
	log "github.com/ligato/cn-infra/logging/logrus"
	"github.com/ligato/vpp-sample-service/bgptol3plugin/writer"
	"github.com/ligato/vpp-sample-service/bgptol3plugin/writer/l3writer"
)

// PluginID of BGP-to-L3 plugin
const PluginID core.PluginName = "bgp-to-l3-plugin"

//TODO use logDep for logging inside plugin
// Plugin with BGP functionality (VPP Agent plugin that servers as BGP-VPP Agent)
// it handles information coming for BGP-Agent channel and sends them transformed to L3 default plugin.
type pluginImpl struct {
	inChan          chan bgp.ReachableIPRoute
	loopStoppedChan chan struct{}
	writer          writer.Writer
}

// New creates Plugin with BGP functionality
func New(bgpChannel chan bgp.ReachableIPRoute) core.NamedPlugin {
	return NewInjectable(bgpChannel, l3writer.New(PluginID))
}

// NewInjectable creates Plugin with BGP functionality with an specific writer implementation
func NewInjectable(bgpChannel chan bgp.ReachableIPRoute, writer writer.Writer) core.NamedPlugin {
	pluginVar := &pluginImpl{
		inChan:          bgpChannel,
		loopStoppedChan: make(chan struct{}),
		writer:          writer,
	}
	return core.NamedPlugin{
		PluginName: PluginID,
		Plugin:     pluginVar,
	}
}

// Init logs attempt of plugin initialization to be sure that plugin is properly recognized. No initialization of plugin is not needed yet.
func (plugin *pluginImpl) Init() error {
	log.DefaultLogger().Info("Initialization of the BGP plugin has completed")
	return nil
}

// AfterInit initializes things depending on proper initialization of L3 default plugin.
func (plugin *pluginImpl) AfterInit() error {
	// prepare VPP for prefix/next hop information sending
	err := plugin.writer.PrepareVPP()
	if err != nil {
		log.DefaultLogger().Errorf("Failed to apply initial VPP configuration: %v", err)
	} else {
		log.DefaultLogger().Info("Successfully applied initial VPP configuration")
	}

	// handle information from BGP-Agent
	go plugin.run()

	log.DefaultLogger().Info("AfterInit of the BGP plugin has completed")
	return nil
}

// Close logs attempt to close this plugin. No resources need cleanup yet.
func (plugin *pluginImpl) Close() error {
	plugin.stop()
	log.DefaultLogger().Info("Closed BGP plugin (in BGP-VPP-Agent)")
	return nil
}

// run is plugin endless loop that handles BGP data from channel.
func (plugin *pluginImpl) run() {
	// notify close that loop is stopped
	defer func() {
		plugin.loopStoppedChan <- struct{}{}
	}()

	// start handler loop
	log.DefaultLogger().Debug("Starting incoming BGP information handler(in BGP-VPP-Agent)")
	for {
		info, openChannel := <-plugin.inChan
		if !openChannel {
			log.DefaultLogger().Debug("Stopping incoming BGP information handler(in BGP-VPP-Agent)")
			return //handler should end its run
		}
		log.DefaultLogger().Debugf("BGP information handler(in BGP-VPP-Agent) received new data: %v", info)
		plugin.sendInformationToVPP(&info)
	}
}

// stop stops plugin loop (and therefore go routine that consists only from this plugin loop)
// This function is blocking until plugin loop stops operating.
func (plugin *pluginImpl) stop() {
	close(plugin.inChan)
	<-plugin.loopStoppedChan //wait for loop to stop
}

// sendInformationToVPP send BGP information from BGP-Agent to VPP.
func (plugin *pluginImpl) sendInformationToVPP(info *bgp.ReachableIPRoute) {
	err := plugin.writer.SendStaticRouteToVPP(info)
	if err != nil {
		log.DefaultLogger().Errorf("Failed to apply BGP information to VPP configuration: %v", err)
	} else {
		log.DefaultLogger().Debug("Successfully applied BGP information to VPP configuration")
	}
}
