package dagger

import "net"

// Packet is an interface that serialize the transfer data
type Packet interface {
	Serialize() []byte
}

// Protocol is an interface that reads the packet data
type Protocol interface {
	ReadPacket(conn *net.TCPConn) (Packet, error)
}
