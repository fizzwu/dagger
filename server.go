package dagger

import "sync"
import "net"
import "time"

// Server ...
type Server struct {
	protocol  Protocol        // custom packet protocol
	callback  SessionCallback // custom server callback
	exitChan  chan struct{}   // notify all goroutines to close
	waitGroup *sync.WaitGroup // wait for all goroutines
}

// NewServer inits a new server instance
func NewServer(protocol Protocol) *Server {
	return &Server{
		protocol:  protocol,
		exitChan:  make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
	}
}

// Serve is a Server run method
func (server *Server) Serve(listener *net.TCPListener, timeout time.Duration) {
	server.waitGroup.Add(1)
	defer func() {
		listener.Close()
		server.waitGroup.Done()
	}()

	for {
		select {
		case <-server.exitChan:
			return
		default:
		}

		listener.SetDeadline(time.Now().Add(timeout))

		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		server.waitGroup.Add(1)
		go func() {
			newSession(conn, server).Work()
			server.waitGroup.Done()
		}()
	}
}
