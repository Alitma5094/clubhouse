package ws

import "sync"

type Manager struct {
	Clients ClientList
	// Using a syncMutex here to be able to lock state before editing clients
	sync.RWMutex
	Handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		Clients:  make(ClientList),
		Handlers: make(map[string]EventHandler),
	}
	return m
}
