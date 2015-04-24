package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pcx/st-agent/conf"
	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/server"
)

func main() {
	MachineID := flag.String("MachineID", "", "Machine ID provided by web dashboard")
	AuthToken := flag.String("AuthToken", "", "Auth token provided by web dashboard")
	HubURL := flag.String("HubURL", "", "Hub url to connect to, format: scheme://host:port")

	flag.Parse()

	config, err := conf.GetConfig(*MachineID, *AuthToken, *HubURL)
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

	s := server.NewServer(config)
	s.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	log.Infof("Gracefully shutting down")
	s.Stop()
}
