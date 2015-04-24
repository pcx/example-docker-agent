package hub

import (
	"net/url"
	"time"

	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
	"github.com/pcx/st-agent/pkg"
)

type Hub struct {
	URL *url.URL
}

func NewHub(hubURL *url.URL) *Hub {
	// hubURL.Scheme = "http"
	return &Hub{URL: hubURL}
}

type Heartbeat struct {
	Machine   string    `json:"machine"`
	AuthToken string    `json:"auth_token"`
	Timestamp time.Time `json:"timestamp"`
}

func newHeartbeat(m *machine.Machine) *Heartbeat {
	return &Heartbeat{
		Machine:   m.MachineID,
		AuthToken: m.AuthToken,
		Timestamp: time.Now().UTC()}
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

	log.Infof("Heartbeat tick success, response is: %v", body)
	return nil
}
