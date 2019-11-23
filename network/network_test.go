package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

var rps uint64

func TestConnection(t *testing.T) {
	for i := 0; i<50; i++ {
		go func() {
			conn, err := net.Dial("tcp", ":8964")
			if err != nil {
				t.Fatal(err)
			}
			for {
				fmt.Fprintf(conn, "PING\r\n")
				_, err := bufio.NewReader(conn).ReadBytes('\n')
				if err != nil {
					t.Fatal(err)
				}
				atomic.AddUint64(&rps, 1)
			}
		}()
	}
	select {
	case <- time.After(5 * time.Second):
		log.Println(rps/5)
		return
	}
}
