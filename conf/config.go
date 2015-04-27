package conf

import (
	"errors"
	"net/url"

	"github.com/codegangsta/cli"

	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/proxy"
)

type Config struct {
	MachineID string
	AuthToken string
	HubURL    *url.URL
	DMan      *proxy.DockerManager
}

func MakeConfig(ctx *cli.Context) (*Config, error) {
	// TODO: Use values from config file if cli ones are empty
	// return nil & error on error

	MachineID := ctx.String("MachineID")
	HubURLStr := ctx.String("HubURL")
	DockerHost := ctx.String("DockerHost")
	DockerTLSVerify := ctx.Bool("DockerTLSVerify")
	DockerCertPath := ctx.String("DockerCertPath")

	if MachineID == "" || HubURLStr == "" {
		return nil, errors.New("Invalid command usage. -m and -h are required.")
	}

	HubURL, err := url.Parse(HubURLStr)
	if err != nil {
		return nil, err
	}

	if (DockerTLSVerify && DockerCertPath == "") ||
		(DockerCertPath != "" && !DockerTLSVerify) {
		return nil, errors.New("DockerTLSVerify and DockerCertPath need to be set together")
	}

	dMan := proxy.NewDockerManager(DockerHost, DockerTLSVerify, DockerCertPath)
	conf := &Config{
		MachineID: MachineID,
		HubURL:    HubURL,
		DMan:      dMan,
	}
	log.Infof("config set, MachineID=%s, HubURL=%s", MachineID, HubURL.String())

	return conf, nil
}
