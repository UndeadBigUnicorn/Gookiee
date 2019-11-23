package main

import (
	"github.com/UndeadBigUnicorn/Gookiee/network"
	"log"
)

func main() {
	log.Fatal(network.NewDefaultNetworkManager().Serve())
}

//The Lshortfile flag includes file name and line number in log messages.
func init() {
	log.SetFlags(log.Lshortfile)
}
