package server

import (
	"github.com/pcx/st-agent/agent"
)

type Server struct {
	agent *agent.Agent
	conf  *Config

	stop chan bool
}

func NewServer(conf *Config) *Server {
	return &Server{
		conf:  conf,
		agent: agent.NewAgent(),
		stop:  make(chan bool)}
}

func (s *Server) Start() {
	go s.agent.HeartbeatRefresh()
}
