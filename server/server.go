package server

import (
	"github.com/pcx/st-agent/agent"
	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/hub"
	"github.com/pcx/st-agent/machine"
)

type Server struct {
	config  *conf.Config
	agent   *agent.Agent
	machine *machine.Machine
	hbMan   *hub.HeartbeatManager
	stop    chan bool
}

func NewServer(config *conf.Config) *Server {
	return &Server{
		config:  config,
		agent:   agent.NewAgent(),
		machine: machine.NewMachine(config),
		hbMan:   hub.NewHeartbeatManager(config),
		stop:    make(chan bool),
	}
}

func (s *Server) Start() {
	go s.machine.MachineStateRefresh(s.stop)
	go s.hbMan.Beat(s.stop)
}

func (s *Server) Stop() {
	s.stop <- true
	close(s.stop)
}
