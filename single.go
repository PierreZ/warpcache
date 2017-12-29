package warpcache

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
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
	cache.selector = s
	cache.config = c
	cache.Errors = make(chan error)
	cache.cache.done = make(chan bool)

	err = cache.initiate()
	if err != nil {
		return nil, err
	}

	go cache.watch()

	return &cache, nil
}

func (c *SingleCache) initiate() error {

	body, err := generateFetchSingleWarpScript(c.config.ReadToken, c.selector.String())

	if err != nil {
		return err
	}
	var resp *http.Response
	resp, err = c.config.HTTPClient.Post(c.config.HTTPProtocol+"://"+os.Getenv("ENDPOINT")+"/api/v0/exec", "", strings.NewReader(body))
	if err != nil {
		return err
	}

	if resp.StatusCode > 200 {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}
		return errors.New(string(dump))
	}

	defer resp.Body.Close()

	money := make([]float64, 1)

	err = json.NewDecoder(resp.Body).Decode(&money)
	if err != nil {
		return err
	}
	log.Println("money is at", money)
	c.v = money[0]
	return nil
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

		c.mux.Lock()
		c.v = value
		c.mux.Unlock()
	}
	goto beginning
}
