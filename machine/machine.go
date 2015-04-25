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
	stateUpdateInterval = 5 * time.Second
)

type Machine struct {
	MachineID string
	AuthToken string

	dMan  *proxy.DockerManager
	state *MachineState
}

type MachineState struct {
	Containers []docker.APIContainers
	Timestamp  time.Time

	mut sync.RWMutex
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
	mach.state.mut.RLock()
	defer mach.state.mut.RUnlock()

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

	mach.state.mut.Lock()
	defer mach.state.mut.Unlock()

	log.Debug("MachineState update tick triggered")
	mach.state.Containers = mach.dMan.ListContainers(listContainerOpts)
}

func (mach *Machine) MachineStateRefresh(stop chan bool) {
	// trigger first update immediately
	mach.machineStateUpdate()

	updateTicker := time.Tick(stateUpdateInterval)
	for {
		select {
		case <-updateTicker:
			mach.machineStateUpdate()
		case <-stop:
			log.Info("Stopping MachineState update due to stop signal")
		}
	}

}
