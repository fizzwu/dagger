package telnet

import (
	"log"

	"github.com/fizzwu/dagger"
)

type TelnetCallback struct {
}

func (this *TelnetCallback) OnConnect(s *dagger.Session) bool {
	addr := s.RawConn().RemoteAddr()
	log.Println("client connected:", addr)
	s.SendPacket(NewTelnetPacket(UndefinedCmd, []byte("Welcome!")), 0)
	return true
}

func (this *TelnetCallback) OnMessage(s *dagger.Session, p dagger.Packet) bool {
	packet := p.(*TelnetPacket)
	flag := packet.Flag()
	command := packet.Data()
	log.Println("flag:", flag)

	switch flag {
	case EchoCmd:
		s.SendPacket(NewTelnetPacket(EchoCmd, command), 0)
	default:
		s.SendPacket(NewTelnetPacket(UndefinedCmd, []byte("Unknown command")), 0)
	}
	return true
}

func (this *TelnetCallback) OnClose(s *dagger.Session) {
	log.Println("Someone leaves..")
}
