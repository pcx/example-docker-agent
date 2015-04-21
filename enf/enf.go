package main

import (
	"flag"

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
		flag.PrintDefaults()
		log.Fatal("exiting")
	}
	s := server.NewServer(conf)
	s.Start()
}
