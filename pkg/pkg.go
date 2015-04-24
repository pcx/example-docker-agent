package pkg

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JSONRequest(url string, content interface{}) (resp *http.Response, body string, err error) {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contentJSON))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, "", err
	}

	bodyInBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	body = string(bodyInBytes)

	return
}
