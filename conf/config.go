package conf

import (
	"fmt"
	"net/url"

	"github.com/pcx/st-agent/log"
)

type Config struct {
	MachineID string
	AuthToken string
	HubURL    *url.URL
}

func GetConfig(MachineID string, AuthToken string, HubURLStr string) (*Config, error) {
	// TODO: Use values from config file if cli ones are empty
	// return nil & error on error

	if MachineID == "" || AuthToken == "" || HubURLStr == "" {
		return nil, fmt.Errorf("Must provide flags MachineID, AuthToken, HubURL")
	}

	HubURL, err := url.Parse(HubURLStr)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		MachineID: MachineID,
		AuthToken: AuthToken,
		HubURL:    HubURL}

	log.Infof("config set, MachineID=%s, AuthToken=%s, HubURL=%s",
		MachineID, AuthToken, HubURL.String())

	return conf, nil
}
