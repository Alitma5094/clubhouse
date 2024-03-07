package ws

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

var (
	/**
	WebsocketUpgrader is used to upgrade incoming HTTP requests into a persistent websocket connection
	*/
	WebsocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Grab the request origin
			// origin := r.Header.Get("Origin")

			// switch origin {
			// case "http://localhost:8080":
			// 	return true
			// default:
			// 	return false
			// }
			return true
		},
	}
)

type Manager struct {
	Clients ClientList

	// Using a syncMutex here to be able to lock state before editing clients
	// Could also use Channels to block
	sync.RWMutex
	Handlers map[string]EventHandler
}

// NewManager is used to initialize all the values inside the manager
func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		Handlers: make(map[string]EventHandler),
	}
	return m
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	log.Println(event.Type)
	if handler, ok := m.Handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// addClient will add clients to our clientList
func (m *Manager) AddClient(client *Client) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.Clients[client] = true
}

// removeClient will remove the client and clean up
func (m *Manager) RemoveClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.Clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.Clients, client)
	}
}
