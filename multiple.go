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

// MultipleCache is watching other multiples GTS
type MultipleCache struct {
	cache
	pivot string
	v     map[string]float64
}

// Get is returning the latest value for a MultipleCache
func (c *MultipleCache) Get(label string) float64 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v[label]
}

// Set is setting a new value. Cache is updated and a new datapoint is pushed
func (c *MultipleCache) Set(label string, f float64) {
	c.mux.Lock()
	c.v[label] = f
	c.mux.Unlock()

	log.Println("pivot:", c.pivot, ",label:", label)

	// Pushing datapoint
	gts := warp.NewGTS(c.cache.selector.Classname).AddLabel(c.pivot, label)
	gts.AddDatapoint(time.Now(), f)

	err := gts.Push(c.config.HTTPProtocol+"://"+c.config.Endpoint, c.config.WriteToken)
	if err != nil {
		c.Errors <- err
	}
}

// Inc is incrementing the value. Cache is updated and a new datapoint is pushed
func (c *MultipleCache) Inc(label string) {
	f := c.Get(label)
	c.Set(label, f+1)
}

func (c *MultipleCache) close() {
	c.cache.done <- true
}

// NewMultipleCache is creating a new MultipleCache
func NewMultipleCache(s Selector, pivot string, c Configuration) (*MultipleCache, error) {

	c.setDefault()

	cache := MultipleCache{}
	cache.v = map[string]float64{}
	cache.selector = s
	cache.config = c
	cache.Errors = make(chan error)
	cache.cache.done = make(chan bool)
	cache.pivot = pivot

	err := cache.initiate()
	if err != nil {
		return nil, err
	}

	go cache.watch()

	return &cache, nil
}

func (c *MultipleCache) initiate() error {

	body, err := generateFetchMultipleWarpScript(c.config.ReadToken, c.selector.String(), c.pivot)

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

	assetsResponse := make([]map[string]float64, 1)

	err = json.NewDecoder(resp.Body).Decode(&assetsResponse)
	if err != nil {
		return err
	}
	assets := assetsResponse[0]
	for name, number := range assets {
		c.v[name] = number
	}
	return nil
}

func (c *MultipleCache) watch() {

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
		labels := make(map[string]string)
		var value float64
		_, labels, value, err = parseInputFormat(string(message))
		if err != nil {
			c.Errors <- err
			continue
		}
		c.mux.Lock()
		c.v[labels[c.pivot]] = value
		c.mux.Unlock()

	}
	goto beginning
}
