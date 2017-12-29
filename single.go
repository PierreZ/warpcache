package warpcache

import (
	"net/url"
	"time"

	warp "github.com/PierreZ/Warp10Exporter"
	"github.com/gorilla/websocket"
)

// SingleCache is a cache for a single GTS.
// To watch multiple GTS at the same time, please use MultipleCache
type SingleCache struct {
	cache
	v float64
}

// Get is returning the latest value for a SingleCache
func (c *SingleCache) Get() float64 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v
}

// Set is setting a new value. Cache is updated and a new datapoint is pushed
func (c *SingleCache) Set(f float64) {
	c.mux.Lock()
	c.v = f
	c.mux.Unlock()

	// Pushing datapoint
	gts := warp.NewGTS(c.cache.selector.Classname)
	gts.Labels = c.selector.Labels
	gts.AddDatapoint(time.Now(), f)

	err := gts.Push(c.config.HTTPProtocol+"://"+c.config.Endpoint, c.config.WriteToken)
	if err != nil {
		c.Errors <- err
	}
}

// Inc is incrementing the value. Cache is updated and a new datapoint is pushed
func (c *SingleCache) Inc() {
	f := c.Get()
	c.Set(f + 1)
}

func (c *SingleCache) close() {
	c.cache.done <- true
}

// NewSingleCache is creating a new SingleCache
func NewSingleCache(s Selector, c Configuration) (*SingleCache, error) {

	c.setDefault()

	// Checking configuration, and the fact that it's a Single GTS
	err := checkSingleGTS(c, s.String())
	if err != nil {
		return nil, err
	}

	cache := SingleCache{}
	cache.config = c
	cache.Errors = make(chan error)
	cache.cache.done = make(chan bool)

	go cache.watch()

	return &cache, nil
}

func (c *SingleCache) watch() {

beginning:
	var err error

	u := url.URL{Scheme: c.config.WebSocketProtol, Host: c.config.Endpoint, Path: "/api/v0/plasma"}

	c.ws, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		c.Errors <- err
	}

	defer close(c.Errors)
	defer close(c.cache.done)

	err = c.ws.WriteMessage(websocket.TextMessage, []byte("SUBSCRIBE "+c.config.ReadToken+" "+c.selector.String()))
	if err != nil {
		c.Errors <- err
	}

	for {

		// TODO: close websocket
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.Errors <- err
			break
		}
		var value float64
		_, _, value, err = parseInputFormat(string(message))
		if err != nil {
			c.Errors <- err
			continue
		}
		c.Set(value)
	}
	goto beginning
}
