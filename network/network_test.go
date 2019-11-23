package network

import "testing"

func TestOpenConnection(t *testing.T) {
	nm := NewNetworkManager(NetConfig{
		Host: "localhost",
		Port: "8080",
		Type: "tcp",
	})

	nm.OpenConnection()
}
