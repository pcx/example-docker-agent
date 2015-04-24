package heart

import (
	"errors"
	"time"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/hub"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
)

const (
	beatInterval = 10 * time.Second
)

type Manager struct {
	hub     *hub.Hub
	machine *machine.Machine
}

func NewManager(config *conf.Config) *Manager {
	return &Manager{
		hub:     hub.NewHub(config.HubURL),
		machine: machine.NewMachine(config.MachineID, config.AuthToken)}
}

func (m *Manager) Beat(stop chan bool) error {
	// trigger first heartbeat immediately
	m.beat()

	beatChan := time.Tick(beatInterval)
	for {
		select {
		case <-beatChan:
			m.beat()
		case <-stop:
			return errors.New("Stopping heartbeat due to stop signal")
		}
	}
}

func (m *Manager) beat() {
	log.Debug("Hearbeat tick triggered")
	ms := m.machine.GetState()
	err := m.hub.SetMachineState(*ms)
	if err != nil {
		log.Errorf("Failed heatbeat: %v", err)
	}
}
