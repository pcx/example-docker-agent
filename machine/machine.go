package machine

import (
	"github.com/pcx/st-agent/conf"
)

type Machine struct {
	MachineID string
	AuthToken string
}

func NewMachine(config *conf.Config) *Machine {
	return &Machine{
		MachineID: config.MachineID,
		AuthToken: config.AuthToken}
}

func (m *Machine) GetState() *Machine {
	return m
}
