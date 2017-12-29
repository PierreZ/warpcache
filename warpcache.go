package warpcache

import (
	"sync"

	"github.com/gorilla/websocket"
)

type cache struct {
	mux      sync.Mutex
	config   Configuration
	ws       *websocket.Conn
	selector Selector
	Errors   chan error
	done     chan bool
}
