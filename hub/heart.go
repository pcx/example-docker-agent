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
	MachineID    string                `json:"machine_id"`
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
		Timestamp:    time.Now().UTC(),
		MachineState: hbMan.mach.GetState(),
	}
}

func (hbMan *HeartbeatManager) beat() {
	hb := hbMan.newHeartbeat()
	beatURL := hbMan.hub.URL.String() + "/heartbeat/"
	resp, body, err := pkg.JSONRequest(beatURL, hb, false)
	if err != nil {
		log.Errorf("Failed sending heartbeat: %v", err)
	} else if resp.StatusCode != 201 {
		log.Errorf("Unexpected status code: %v, Response is : %v", resp.StatusCode, body)
	} else if resp.Body == nil {
		// just to cover any case where res.Body is still nil
		log.Fatal("Fatal, heartbeat response should be empty. Dying!")
	} else {
		defer resp.Body.Close()
		log.Infof("Heartbeat tick success, response is: %v", body)
	}
}

func (hbMan *HeartbeatManager) Beat(stop chan bool) {
	// trigger first heartbeat immediately
	hbMan.beat()

	beatChan := time.Tick(beatInterval)
	for {
		select {
		case <-beatChan:
			hbMan.beat()
		case <-stop:
			log.Info("Stopping heartbeat due to stop signal")
		}
	}
}
