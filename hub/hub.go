package hub

import (
	"fmt"
	"net/url"

	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
	"github.com/pcx/st-agent/pkg"
)

type Hub struct {
	URL *url.URL
}

func NewHub(hubURL *url.URL) *Hub {
	return &Hub{URL: hubURL}
}

// Machine argument is passed by value, so that SetMachineState
// will have a private copy that cannot be modified by other code
func (h *Hub) SetMachineState(m machine.Machine) error {
	log.Infof("POST to sent")

	hb := newHeartbeat(&m)
	resp, body, err := pkg.JSONRequest(h.URL.String()+"/heartbeat/", hb)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return fmt.Errorf("Unexpected status code: %v, Response is : %v", resp.StatusCode, body)
	}
	log.Infof("Heartbeat tick success, response is: %v", body)
	return nil
}
