package network

import (
	"github.com/tidwall/evio"
	"github.com/tidwall/redcon"
	"log"
	"strings"
)

type NetEvents struct {
	events evio.Events
}

// Create NetEvents based on NetConfig
func (nm *NetManager) useEvents() {

	// set number of loops
	nm.Events.events.NumLoops = nm.Config.Loops

	// This option is only available when nm.Events.NumLoops is set.
	// validate inputted load-balancing
	switch nm.Config.Balance {
	default:
		nm.Events.events.LoadBalance = evio.Random
	case "random":
		nm.Events.events.LoadBalance = evio.Random
	case "round-robin":
		nm.Events.events.LoadBalance = evio.RoundRobin
	case "least-connections":
		nm.Events.events.LoadBalance = evio.LeastConnections
	}

	// what's going on while server starts
	nm.Events.events.Serving = nm.onServe()
	// what's going on while someone opens a new connection
	nm.Events.events.Opened = nm.onOpen()
	// what's going on while someone closes the connection
	nm.Events.events.Closed = nm.onClose()
	// what's going on while someone send data
	nm.Events.events.Data = nm.onData()

}

type serveFunc func(srv evio.Server) (action evio.Action)
type openFunc func(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action)
type closeFunc func(ec evio.Conn, err error) (action evio.Action)
type dataFunc func(ec evio.Conn, in []byte) (out []byte, action evio.Action)

func (nm *NetManager) onServe() serveFunc {
	return func(srv evio.Server) (action evio.Action) {
		log.Printf("gookiee started on port %d (loops: %d)", nm.Config.Port, srv.NumLoops)
		return
	}
}

func (nm *NetManager) onOpen() openFunc {
	return func(ec evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
		// TODO: create own logging system, because it will be tons of logs
		//log.Printf("opened: %v\n", ec.RemoteAddr())
		// TODO: the biggest question is if I really need to store information about connection
		// because all I need is incoming data
		conn := new(NetConnection)
		nm.Connections.AddConnection(ec.RemoteAddr().String(), conn)
		ec.SetContext(conn)
		return
	}
}

func (nm *NetManager) onClose() closeFunc {
	return func(ec evio.Conn, err error) (action evio.Action) {
		log.Printf("closed: %v\n", ec.RemoteAddr())
		nm.Connections.DeleteConnection(ec.RemoteAddr().String())
		return
	}
}

func (nm *NetManager) onData() dataFunc {
	return func(ec evio.Conn, in []byte) (out []byte, action evio.Action) {
		if in == nil {
			log.Printf("wake from %s\n", ec.RemoteAddr())
			return nil, evio.Close
		}
		c := ec.Context().(*NetConnection)
		data := c.is.Begin(in)
		var n int
		var complete bool
		var err error
		var args [][]byte
		for action == evio.None {
			complete, args, _, data, err = redcon.ReadNextCommand(data, args[:0])
			if err != nil {
				action = evio.Close
				out = redcon.AppendError(out, err.Error())
				break
			}
			if !complete {
				break
			}
			if len(args) > 0 {
				n++
				switch strings.ToUpper(string(args[0])) {
				default:
					out = redcon.AppendError(out, "ERR unknown command '"+string(args[0])+"'")
				case "PING":
					if len(args) > 2 {
						out = redcon.AppendError(out, "ERR wrong number of arguments for '"+string(args[0])+"' command")
					} else if len(args) == 2 {
						out = redcon.AppendBulk(out, args[1])
					} else {
						out = redcon.AppendString(out, "PONG")
					}
				case "WAKE":
					go ec.Wake()
					out = redcon.AppendString(out, "OK")
				case "ECHO":
					if len(args) != 2 {
						out = redcon.AppendError(out, "ERR wrong number of arguments for '"+string(args[0])+"' command")
					} else {
						out = redcon.AppendBulk(out, args[1])
					}
				case "SHUTDOWN":
					out = redcon.AppendString(out, "OK")
					action = evio.Shutdown
				case "QUIT":
					out = redcon.AppendString(out, "OK")
					action = evio.Close
				//case "GET":
				//	if len(args) != 2 {
				//		out = redcon.AppendError(out, "ERR wrong number of arguments for '"+string(args[0])+"' command")
				//	} else {
				//		key := string(args[1])
				//		mu.Lock()
				//		val, ok := keys[key]
				//		mu.Unlock()
				//		if !ok {
				//			out = redcon.AppendNull(out)
				//		} else {
				//			out = redcon.AppendBulkString(out, val)
				//		}
				//	}
				//case "SET":
				//	if len(args) != 3 {
				//		out = redcon.AppendError(out, "ERR wrong number of arguments for '"+string(args[0])+"' command")
				//	} else {
				//		key, val := string(args[1]), string(args[2])
				//		mu.Lock()
				//		keys[key] = val
				//		mu.Unlock()
				//		out = redcon.AppendString(out, "OK")
				//	}
				//case "DEL":
				//	if len(args) < 2 {
				//		out = redcon.AppendError(out, "ERR wrong number of arguments for '"+string(args[0])+"' command")
				//	} else {
				//		var n int
				//		mu.Lock()
				//		for i := 1; i < len(args); i++ {
				//			if _, ok := keys[string(args[i])]; ok {
				//				n++
				//				delete(keys, string(args[i]))
				//			}
				//		}
				//		mu.Unlock()
				//		out = redcon.AppendInt(out, int64(n))
				//	}
				//case "FLUSHDB":
				//	mu.Lock()
				//	keys = make(map[string]string)
				//	mu.Unlock()
				//	out = redcon.AppendString(out, "OK")
				}
			}
		}
		c.is.End(data)
		return
	}
}
