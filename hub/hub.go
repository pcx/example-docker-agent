package hub

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pcx/st-agent/log"
	"github.com/pcx/st-agent/machine"
)

type Hub struct {
	url string
}

func NewHub() *Hub {
	return &Hub{url: "http://lvh.me:8000"}
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
	hbJSON, err := json.Marshal(hb)
	if err != nil {
		return err
	}
	log.Debugf("Heartbeat is: %v", string(hbJSON))

	req, err := http.NewRequest("POST", h.url+"/heartbeat/", bytes.NewBuffer(hbJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("Heartbeat req success, response is: %v", string(body))
	return nil
}
