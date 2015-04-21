package server

import (
	"fmt"

	"github.com/pcx/st-agent/log"
)

type Config struct {
	machineId string
	authToken string
}

func GetConfig(machineId string, authToken string) (*Config, error) {
	// TODO: Use values from config file if cli ones are empty
	// return nil & error on error

	if machineId == "" || authToken == "" {
		return nil, fmt.Errorf("either is empty machineid=%s authtoken=%s", machineId, authToken)
	} else {
		log.Infof("config set, machineId %s and authToken %s", machineId, authToken)
	}

	return &Config{machineId: machineId, authToken: authToken}, nil
}
