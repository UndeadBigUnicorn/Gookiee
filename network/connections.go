package network

import (
	"github.com/tidwall/evio"
	"sync"
)

type NetConnection struct {
	is   evio.InputStream
	addr string
}

type Connections struct {
	conns ConnectionsMap
	mutex sync.RWMutex
}

type ConnectionsMap map[string] *NetConnection

func newDefaultConnections() *Connections {
	return &Connections{
		conns: make(ConnectionsMap),
		mutex: sync.RWMutex{},
	}
}

func (c *Connections) AddConnection(addr string, conn *NetConnection) {
	c.mutex.Lock()
	c.conns[addr] = conn
	c.mutex.Unlock()
}

func (c *Connections) DeleteConnection(addr string) {
	c.mutex.RLock()
	delete(c.conns, addr)
	c.mutex.RUnlock()
}

func (c *Connections) LoadConnection(addr string) *NetConnection {
	return c.conns[addr]
}
