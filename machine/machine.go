package machine

type Machine struct {
	MachineID string
	AuthToken string
}

func NewMachine(machineID string, authToken string) *Machine {
	return &Machine{
		MachineID: machineID,
		AuthToken: authToken}
}

func (m *Machine) GetState() *Machine {
	return m
}
