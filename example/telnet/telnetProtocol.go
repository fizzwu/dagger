package telnet

import (
	"net"
	"strings"

	"bytes"

	"github.com/fizzwu/dagger"
)

var endTag = []byte("\r\n") // telnet command end tag

const (
	// UndefinedCmd ..
	UndefinedCmd = "undefined"
	// ExitCmd ..
	ExitCmd = "exit"
	// EchoCmd ...
	EchoCmd = "echo"
)

// TelnetPacket ..
type TelnetPacket struct {
	flag string
	data []byte
}

// Flag is the flag getter
func (tp *TelnetPacket) Flag() string {
	return tp.flag
}

// Data is the data getter
func (tp *TelnetPacket) Data() []byte {
	return tp.data
}

// Serialize ..
func (tp *TelnetPacket) Serialize() []byte {
	buf := tp.data
	buf = append(buf, endTag...)
	return buf
}

// NewTelnetPacket inits a TelnetPacket instance
func NewTelnetPacket(flag string, data []byte) *TelnetPacket {
	return &TelnetPacket{
		flag: flag,
		data: data,
	}
}

// TelnetProtocol ..
type TelnetProtocol struct {
}

// ReadPacket ..
func (tp *TelnetProtocol) ReadPacket(conn *net.TCPConn) (dagger.Packet, error) {
	fullBuf := bytes.NewBuffer([]byte{})
	for {
		// read from the conn and put the data into a buffer until endTag presents
		data := make([]byte, 1024)
		connReadLenth, err := conn.Read(data)
		if err != nil {
			return nil, err
		}

		if connReadLenth == 0 {
			return nil, dagger.ErrConnClosed
		}

		fullBuf.Write(data[:connReadLenth])
		index := bytes.Index(fullBuf.Bytes(), endTag)
		if index > -1 {
			inputCommand := fullBuf.Next(index) // get command bytes
			fullBuf.Next(2)                     // skip endTag
			commandList := strings.Split(string(inputCommand), " ")
			if len(commandList) == 2 {
				return NewTelnetPacket(commandList[0], []byte(commandList[1])), nil
			}
			if len(commandList) == 1 && commandList[0] == "quit" {
				return NewTelnetPacket(ExitCmd, []byte{}), nil
			}
			if len(commandList) == 1 && commandList[0] == "echo" {
				return NewTelnetPacket(EchoCmd, []byte{}), nil
			}
			return NewTelnetPacket(UndefinedCmd, inputCommand), nil
		}
	}
}
