package dagger

import (
	"net"
	"sync"
	"sync/atomic"
)

// Session ...
type Session struct {
	server    *Server       // server pointer
	conn      *net.TCPConn  // raw conn
	closeChan chan struct{} // close channel
	closeFlag int32         // close flag, set to 1 when closed
	closeOnce sync.Once     // close session only once
	sendChan  chan Packet   // packet send chan
	recvChan  chan Packet   // packet receive chan
}

func newSession(conn *net.TCPConn, server *Server) *Session {
	return &Session{
		server:    server,
		conn:      conn,
		closeChan: make(chan struct{}),
		sendChan:  make(chan Packet),
		recvChan:  make(chan Packet),
	}
}

// RawConn is the conn instance getter
func (s *Session) RawConn() *net.TCPConn {
	return s.conn
}

// IsClosed indicates whether session is closed
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

// Close closes the session
func (s *Session) Close() {
	atomic.StoreInt32(&s.closeFlag, 1)
	close(s.closeChan)
	close(s.sendChan)
	close(s.recvChan)
	s.conn.Close()
	s.server.callback.OnClose(s)
}

// Work is the session handler
func (s *Session) Work() {
	if s.server.callback.OnConnect(s) == false {
		return
	}

	parallel(s.handleLoop, s.server.waitGroup)
	parallel(s.readLoop, s.server.waitGroup)
	parallel(s.writeLoop, s.server.waitGroup)
}

// read packet from the conn and send to recvChan
func (s *Session) readLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		case <-s.server.exitChan:
			return
		case <-s.closeChan:
			return
		default:
		}

		p, err := s.server.protocol.ReadPacket(s.conn)
		if err != nil {
			return
		}
		s.recvChan <- p
	}
}

// receive data from the recvChan, write to conn
func (s *Session) writeLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		case <-s.server.exitChan:
			return
		case <-s.closeChan:
			return
		case p := <-s.recvChan:
			if s.IsClosed() {
				return
			}
			if _, err := s.conn.Write(p.Serialize()); err != nil {
				return
			}
		}
	}
}

func (s *Session) handleLoop() {
	defer func() {
		recover()
		s.Close()
	}()

	for {
		select {
		case <-s.server.exitChan:
			return
		case <-s.closeChan:
			return
		case p := <-s.recvChan:
			if s.IsClosed() {
				return
			}
			if !s.server.callback.OnMessage(s, p) { // session closed
				return
			}
		}
	}
}

// run function in a goroutine with waitgroup managed
func parallel(f func(), wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}
