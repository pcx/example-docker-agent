package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/server"
)

func main() {
	machineId := flag.String("machineId", "", "Machine ID provided by web dashboard")
	authToken := flag.String("authToken", "", "Auth token provided by web dashboard")
	flag.Parse()

	conf, err := server.GetConfig(*machineId, *authToken)
	if err != nil {
		log.Errorf("Unable to parse config: %v", err)
		fmt.Println("Usage: enf [OPTIONS]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	log.EnableTimestamps()
	log.EnableDebug()

	log.Debug("Starting enforcer daemon")

	s := server.NewServer(conf)
	s.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	log.Infof("Gracefully shutting down")
	s.Stop()
}
