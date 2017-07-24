package dagger

// SessionCallback is an interface that describes a server session callback bahaviours
type SessionCallback interface {
	// OnConnect is called when connection is accepted, return false if session closed
	OnConnect(*Session) bool

	// OnMessage is called when connection receives a packet, return false if session closed
	OnMessage(*Session, Packet) bool

	// OnClose is called when the session is closed
	OnClose(*Session)
}
