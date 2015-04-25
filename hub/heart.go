package hub

import (
	"time"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
	"github.com/pcx/st-agent/pkg"
)

const (
	beatInterval = 10 * time.Second
)

type Heartbeat struct {
	MachineID    string                `json:"machine"`
	AuthToken    string                `json:"auth_token"`
	Timestamp    time.Time             `json:"timestamp"`
	MachineState *machine.MachineState `json:"machine_state"`
}

type HeartbeatManager struct {
	hub  *Hub
	mach *machine.Machine
}

func NewHeartbeatManager(config *conf.Config) *HeartbeatManager {
	return &HeartbeatManager{
		hub:  NewHub(config.HubURL),
		mach: machine.NewMachine(config),
	}
}

func (hbMan *HeartbeatManager) newHeartbeat() *Heartbeat {
	return &Heartbeat{
		MachineID:    hbMan.mach.MachineID,
		AuthToken:    hbMan.mach.AuthToken,
		Timestamp:    time.Now().UTC(),
		MachineState: hbMan.mach.GetState(),
	}
}

func (hbMan *HeartbeatManager) beat() {
	hb := hbMan.newHeartbeat()
	beatURL := hbMan.hub.URL.String() + "/heartbeat/"
	resp, body, err := pkg.JSONRequest(beatURL, hb, true)
	if err != nil {
		log.Errorf("Failed sending heartbeat: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		log.Errorf("Unexpected status code: %v, Response is : %v", resp.StatusCode, body)
	}
	log.Infof("Heartbeat tick success, response is: %v", body)
}

func (hbMan *HeartbeatManager) Beat(stop chan bool) {
	// trigger first heartbeat immediately
	hbMan.beat()

	beatChan := time.Tick(beatInterval)
	for {
		select {
		case <-beatChan:
			log.Debug("Hearbeat tick triggered")
			hbMan.beat()
		case <-stop:
			log.Info("Stopping heartbeat due to stop signal")
		}
	}
}
