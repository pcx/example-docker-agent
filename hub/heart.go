package hub

import (
	"time"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
)

const (
	beatInterval = 10 * time.Second
)

type Heartbeat struct {
	Machine   string    `json:"machine"`
	AuthToken string    `json:"auth_token"`
	Timestamp time.Time `json:"timestamp"`
}

func newHeartbeat(m *machine.Machine) *Heartbeat {
	return &Heartbeat{
		Machine:   m.MachineID,
		AuthToken: m.AuthToken,
		Timestamp: time.Now().UTC()}
}

type HeartbeatManager struct {
	hub     *Hub
	machine *machine.Machine
}

func NewHeartbeatManager(config *conf.Config) *HeartbeatManager {
	return &HeartbeatManager{
		hub:     NewHub(config.HubURL),
		machine: machine.NewMachine(config)}
}

func (m *HeartbeatManager) Beat(stop chan bool) {
	// trigger first heartbeat immediately
	m.beat()

	beatChan := time.Tick(beatInterval)
	for {
		select {
		case <-beatChan:
			log.Debug("Hearbeat tick triggered")
			m.beat()
		case <-stop:
			log.Info("Stopping heartbeat due to stop signal")
		}
	}
}

func (m *HeartbeatManager) beat() {
	ms := m.machine
	err := m.hub.SetMachineState(*ms)
	if err != nil {
		log.Errorf("Failed sending heartbeat: %v", err)
	}
}
