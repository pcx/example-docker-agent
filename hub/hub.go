package hub

import (
	"net/url"
)

type Hub struct {
	URL *url.URL
}

func NewHub(hubURL *url.URL) *Hub {
	return &Hub{URL: hubURL}
}
