package machine

import (
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/proxy"
)

const (
	stateUpdateInterval = time.Duration(10) * time.Second
)

type Machine struct {
	MachineID string
	AuthToken string

	mut   sync.RWMutex
	dMan  *proxy.DockerManager
	state *MachineState
}

type MachineState struct {
	Containers []docker.APIContainers `json:"containers"`
	Timestamp  time.Time              `json:"timestamp"`
}

func NewMachine(config *conf.Config) *Machine {
	listContainerOpts := docker.ListContainersOptions{
		All:     true,
		Size:    true,
		Limit:   0,
		Since:   "",
		Before:  "",
		Filters: make(map[string][]string),
	}
	containers := config.DMan.ListContainers(listContainerOpts)
	return &Machine{
		MachineID: config.MachineID,
		AuthToken: config.AuthToken,
		dMan:      config.DMan,
		state: &MachineState{
			Containers: containers,
			Timestamp:  time.Now().UTC(),
		},
	}
}

func (mach *Machine) GetState() *MachineState {
	mach.mut.RLock()
	defer mach.mut.RUnlock()

	return mach.state
}

func (mach *Machine) machineStateUpdate() {
	listContainerOpts := docker.ListContainersOptions{
		All:     true,
		Size:    true,
		Limit:   0,
		Since:   "",
		Before:  "",
		Filters: make(map[string][]string),
	}

	mach.mut.Lock()
	defer mach.mut.Unlock()

	// For timestamp to be as accurate as possible, fetch container list
	// when mutex is locked, though otherwise it isn't necessary
	containers := mach.dMan.ListContainers(listContainerOpts)

	mach.state = &MachineState{
		Containers: containers,
		Timestamp:  time.Now().UTC(),
	}
}

func (mach *Machine) MachineStateRefresh(stop chan bool) {
	// trigger first update immediately
	mach.machineStateUpdate()

	updateTicker := time.Tick(stateUpdateInterval)
	for {
		select {
		case <-updateTicker:
			log.Debugf("MachineState update tick triggered, %v", time.Now().UTC())
			mach.machineStateUpdate()
		case <-stop:
			log.Info("Stopping MachineState update due to stop signal")
		}
	}

}
