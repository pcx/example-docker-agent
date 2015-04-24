package conf

import (
	"errors"
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
	DockerHost := ctx.String("DockerHost")
	DockerTLSVerify := ctx.Bool("DockerTLSVerify")
	DockerCertPath := ctx.String("DockerCertPath")

	if MachineID == "" || AuthToken == "" || HubURLStr == "" {
		return nil, errors.New("Invalid command usage. -m and -a are required.")
	}

	HubURL, err := url.Parse(HubURLStr)
	if err != nil {
		return nil, err
	}

	if (DockerTLSVerify && DockerCertPath == "") ||
		(DockerCertPath != "" && !DockerTLSVerify) {
		return nil, errors.New("DockerTLSVerify and DockerCertPath need to be set together")
	}

	// TODO: Move docker connection creation to separate methid
	// provide ability to recreate docker conn if failed during run
	// program should try 10 times for docker client connection

	// to stop code linters from bothering about existence of 'client'
	client := &docker.Client{}

	if DockerTLSVerify {
		cert := fmt.Sprintf("%s/cert.pem", DockerCertPath)
		key := fmt.Sprintf("%s/key.pem", DockerCertPath)
		ca := fmt.Sprintf("%s/ca.pem", DockerCertPath)
		client, err = docker.NewTLSClient(DockerHost, cert, key, ca)
	} else {
		client, err = docker.NewClient(DockerHost)
	}
	if err != nil {
		return nil, fmt.Errorf("Could not connect to Docker daemon at %s: %v",
			DockerHost, err)
	}
	err = client.Ping()
	if err != nil {
		return nil, err
	} else {
		log.Infof("Successfully established connection to Docker host at: %s", DockerHost)
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
