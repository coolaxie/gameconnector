package gate

import (
	"github.com/coolaxie/gameconnector/log"
	"github.com/coolaxie/gameconnector/network/tcp"
	"errors"
	"io"
)

type TCPAgent struct {
	conn *tcp.TCPConn
}

func newTCPAgent(conn *tcp.TCPConn) tcp.Agent {
	a := new(TCPAgent)
	a.conn = conn

	return a
}

func (a *TCPAgent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			if err == io.EOF {
				log.Release("client(%v) closed", a.conn.RemoteAddr())
			} else {
				log.Error("client(%v) read error(%v)", a.conn.RemoteAddr(), err)
			}
			break
		}

		//TODO transfer message
		log.Release("received message(%v) from %v", string(data), a.conn.RemoteAddr())
	}
}