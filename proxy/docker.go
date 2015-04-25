package proxy

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"

	"github.com/pcx/st-agent/log"
)

type DockerManager struct {
	client   *docker.Client
	isTLS    bool
	certPath string
	host     string
}

func NewDockerManager(host string, isTLS bool, certPath string) *DockerManager {
	var client *docker.Client
	var err error
	if isTLS {
		cert := fmt.Sprintf("%s/cert.pem", certPath)
		key := fmt.Sprintf("%s/key.pem", certPath)
		ca := fmt.Sprintf("%s/ca.pem", certPath)
		client, err = docker.NewTLSClient(host, cert, key, ca)
	} else {
		client, err = docker.NewClient(host)
	}
	if err != nil {
		log.Fatalf("Fatal error connecting to Docker host at %s: %v", host, err)
	}

	return &DockerManager{
		client:   client,
		host:     host,
		isTLS:    isTLS,
		certPath: certPath,
	}
}

func (dMan *DockerManager) tryPing() bool {
	for i := 1; i <= 10; i++ {
		log.Infof("Trying to ping Docker host at %s", dMan.host)
		err := dMan.client.Ping()
		if err == nil {
			log.Infof("Successfully pinged Docker host at: %s", dMan.host)
			return true
		}
		log.Errorf("Failed to ping docker host at %s", dMan.host)
	}
	log.Fatalf("Maximum retries reached while trying to ping Docker host at %s"+
		", program is exiting now", dMan.host)
	return false
}

func (dMan *DockerManager) ListContainers(opts docker.ListContainersOptions) []docker.APIContainers {
	for i := 1; i <= 3; i++ {
		containers, err := dMan.client.ListContainers(opts)
		if err != nil {
			log.Errorf("%s", err)
			if dMan.tryPing() {
				continue
			}
		} else {
			return containers
		}
	}
	log.Fatalf("Maximum retries reached while communicating Docker host at %s"+
		", program is exiting now", dMan.host)

	return nil
}
