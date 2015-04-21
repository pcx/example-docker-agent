package agent

type Agent struct {
}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) HeartbeatRefresh() {
	return
}
