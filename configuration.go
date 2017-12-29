package warpcache

import (
	"net/http"
	"time"
)

// Configuration is holding the configuration for one cache
type Configuration struct {
	ReadToken         string
	WriteToken        string
	Endpoint          string
	WebSocketProtol   string
	HTTPProtocol      string
	ForceSyncInterval *time.Duration
	HTTPClient        *http.Client
}

func (c *Configuration) setDefault() {
	if c.WebSocketProtol == "" {
		c.WebSocketProtol = "wss"
	}

	if c.HTTPProtocol == "" {
		c.HTTPProtocol = "https"
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

}
