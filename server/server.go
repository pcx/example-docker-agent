package server

import (
	"github.com/pcx/st-agent/agent"
	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/heart"
)

type Server struct {
	agent  *agent.Agent
	config *conf.Config
	heart  *heart.Manager

	stop chan bool
}

func NewServer(config *conf.Config) *Server {
	return &Server{
		config: config,
		agent:  agent.NewAgent(),
		heart:  heart.NewManager(config),
		stop:   make(chan bool)}
}

func (s *Server) Start() {
	go s.heart.Beat(s.stop)
}

func (s *Server) Stop() {
	s.stop <- true
	close(s.stop)
}
