package server

import (
	"github.com/pcx/st-agent/agent"
	"github.com/pcx/st-agent/heart"
	"github.com/pcx/st-agent/machine"
)

type Server struct {
	agent *agent.Agent
	conf  *Config
	heart *heart.Manager

	stop chan bool
}

func NewServer(conf *Config) *Server {
	return &Server{
		conf:  conf,
		agent: agent.NewAgent(),
		heart: heart.NewManager(
			machine.NewMachine(
				conf.machineId,
				conf.authToken)),
		stop: make(chan bool)}
}

func (s *Server) Start() {
	go s.heart.Beat(s.stop)
}

func (s *Server) Stop() {
	s.stop <- true
	close(s.stop)
}
