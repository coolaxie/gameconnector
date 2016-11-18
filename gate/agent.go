package gate

import (
	"github.com/coolaxie/gameconnector/log"
	"github.com/coolaxie/gameconnector/network/tcp"
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
			log.Error("read error(%v)", err)
			break
		}

		//TODO transfer message
		log.Release("received message(%v)", string(data))
	}
}