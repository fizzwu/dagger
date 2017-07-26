package dagger

import "sync"
import "net"
import "time"

// Server ...
type Server struct {
	protocol       Protocol        // custom packet protocol
	callback       SessionCallback // custom server callback
	exitChan       chan struct{}   // notify all goroutines to close
	waitGroup      *sync.WaitGroup // wait for all goroutines
	packetSendSize uint32          // session's packet send channel size
	packetRecvSize uint32          // session's packet receive channel size
}

// NewServer inits a new server instance
func NewServer(callback SessionCallback, protocol Protocol, packetSendSize uint32, packetRecvSize uint32) *Server {
	return &Server{
		protocol:       protocol,
		callback:       callback,
		exitChan:       make(chan struct{}),
		waitGroup:      &sync.WaitGroup{},
		packetSendSize: packetSendSize,
		packetRecvSize: packetRecvSize,
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
