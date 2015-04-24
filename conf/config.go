package conf

import (
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"

	"github.com/pcx/st-agent/log"
)

type Config struct {
	MachineID  string
	AuthToken  string
	HubURL     *url.URL
	DockerConn *docker.Client
}

func MakeConfig(ctx *cli.Context) (*Config, error) {
	// TODO: Use values from config file if cli ones are empty
	// return nil & error on error

	MachineID := ctx.String("MachineID")
	AuthToken := ctx.String("AuthToken")
	HubURLStr := ctx.String("HubURL")

	if MachineID == "" || AuthToken == "" || HubURLStr == "" {
		return nil, fmt.Errorf("The flags MachineID, AuthToken & HubURL are required")
	}

	HubURL, err := url.Parse(HubURLStr)
	if err != nil {
		return nil, err
	}

	// TODO: Move docker connection creation to separate methid
	// provide ability to recreate docker conn if failed during run
	// program should try 10 times for docker client connection
	dockerSock := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(dockerSock)
	if err != nil {
		return nil, fmt.Errorf("Could not create Docker connection at %s: %v", dockerSock, err)
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	}

	conf := &Config{
		MachineID:  MachineID,
		AuthToken:  AuthToken,
		HubURL:     HubURL,
		DockerConn: client,
	}

	log.Infof("config set, MachineID=%s, AuthToken=%s, HubURL=%s",
		MachineID, AuthToken, HubURL.String())

	return conf, nil
}
