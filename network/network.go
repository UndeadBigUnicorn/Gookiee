package network

import (
	"fmt"
	"github.com/tidwall/evio"
)

type NetManager struct {
	Config      NetConfig
	Events      NetEvents
	Connections *Connections
}

type NetConfig struct {
	Host       string // which host to use
	Port       int    // server port
	ReusePort  bool   // "reuseport (SO_REUSEPORT)"
	Type       string // TCP/UDP
	UnixSocket string
	Stdlib     bool   // use stdlib

	// The events.NumLoops options sets the number of loops to use for the server.
	// A value greater than 1 will effectively make the server multithreaded for multi-core machines.
	// Which means you must take care when synchronizing memory between event callbacks.
	// Setting to 0 or 1 will run the server as single-threaded.
	// Setting to -1 will automatically assign this value equal to runtime.NumProcs().
	Loops int // num of loops

	// This option is only available when events.NumLoops is set.
	// * Random requests that connections are randomly distributed.
	// * RoundRobin requests that connections are distributed to a loop in a round-robin fashion.
	// * LeastConnections assigns the next accepted connection to the loop with the least number of active connections.
	Balance string // load balancing  default: "random", possible variants: "random, round-robin, least-connections"
}

// Create network manager with default settings
func NewDefaultNetworkManager() *NetManager {
	manager := &NetManager {
		Config:      NewDefaultConfig(),
		Connections: newDefaultConnections(),
	}
	manager.useEvents()
	return manager
}

// Create config with default values
func NewDefaultConfig() NetConfig {
	return NetConfig {
		Host:       "localhost",
		Port:       8964,
		ReusePort:  false,
		Type:       "TCP",
		UnixSocket: "socket",
		Stdlib:     false,
		Loops:      -1,
		Balance:    "random",
	}
}

func NewNetworkManager(config NetConfig) *NetManager {
	manager := &NetManager {
		Config: config,
		Connections: newDefaultConnections(),
	}
	manager.useEvents()
	return manager
}

func (nm *NetManager) Serve() error {
	addrs := []string{fmt.Sprintf("tcp://:%d?reuseport=%t", nm.Config.Port, nm.Config.ReusePort)}
	if nm.Config.UnixSocket != "" {
		addrs = append(addrs, fmt.Sprintf("unix://%s", nm.Config.UnixSocket))
	}
	return evio.Serve(nm.Events.events, addrs...)
}
