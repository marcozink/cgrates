/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package services

import (
	"fmt"
	"sync"

	"github.com/cgrates/cgrates/agents"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/servmanager"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/rpcclient"
)

// NewDNSAgent returns the DNS Agent
func NewDNSAgent(cfg *config.CGRConfig, filterSChan chan *engine.FilterS,
	exitChan chan bool, connMgr *engine.ConnManager) servmanager.Service {
	return &DNSAgent{
		cfg:         cfg,
		filterSChan: filterSChan,
		exitChan:    exitChan,
		connMgr:     connMgr,
	}
}

// DNSAgent implements Agent interface
type DNSAgent struct {
	sync.RWMutex
	cfg         *config.CGRConfig
	filterSChan chan *engine.FilterS
	exitChan    chan bool

	dns     *agents.DNSAgent
	connMgr *engine.ConnManager

	oldListen string
}

// Start should handle the sercive start
func (dns *DNSAgent) Start() (err error) {
	if dns.IsRunning() {
		return fmt.Errorf("service aleady running")
	}

	filterS := <-dns.filterSChan
	dns.filterSChan <- filterS

	dns.Lock()
	defer dns.Unlock()
	dns.oldListen = dns.cfg.DNSAgentCfg().Listen
	dns.dns, err = agents.NewDNSAgent(dns.cfg, filterS, dns.connMgr)
	if err != nil {
		utils.Logger.Err(fmt.Sprintf("<%s> error: <%s>", utils.DNSAgent, err.Error()))
		return
	}
	go func() {
		if err = dns.dns.ListenAndServe(); err != nil {
			utils.Logger.Err(fmt.Sprintf("<%s> error: <%s>", utils.DNSAgent, err.Error()))
			dns.exitChan <- true // stop the engine here
		}
	}()
	return
}

// GetIntenternalChan returns the internal connection chanel
// no chanel for DNSAgent
func (dns *DNSAgent) GetIntenternalChan() (conn chan rpcclient.ClientConnector) {
	return nil
}

// Reload handles the change of config
func (dns *DNSAgent) Reload() (err error) {
	if dns.oldListen == dns.cfg.DNSAgentCfg().Listen {
		return
	}
	if err = dns.Shutdown(); err != nil {
		return
	}
	dns.Lock()
	dns.oldListen = dns.cfg.DNSAgentCfg().Listen
	defer dns.Unlock()
	if err = dns.dns.Reload(); err != nil {
		return
	}
	go func() {
		if err := dns.dns.ListenAndServe(); err != nil {
			utils.Logger.Err(fmt.Sprintf("<%s> error: <%s>", utils.DNSAgent, err.Error()))
			dns.exitChan <- true // stop the engine here
		}
	}()
	return
}

// Shutdown stops the service
func (dns *DNSAgent) Shutdown() (err error) {
	dns.Lock()
	defer dns.Unlock()
	if err = dns.dns.Shutdown(); err != nil {
		return
	}
	dns.dns = nil
	return
}

// IsRunning returns if the service is running
func (dns *DNSAgent) IsRunning() bool {
	dns.RLock()
	defer dns.RUnlock()
	return dns != nil && dns.dns != nil
}

// ServiceName returns the service name
func (dns *DNSAgent) ServiceName() string {
	return utils.DNSAgent
}

// ShouldRun returns if the service should be running
func (dns *DNSAgent) ShouldRun() bool {
	return dns.cfg.DNSAgentCfg().Enabled
}
