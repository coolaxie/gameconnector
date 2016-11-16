package tcp

import (
	"github.com/coolaxie/gameconnector/log"
	"net"
)

type Agent struct {
	conn      net.Conn
	writeChan chan []byte
}

func newAgent(conn net.Conn) *Agent {
	a := new(Agent)
	a.conn = conn
	a.writeChan = make(chan []byte, 2000)

	go a.writeLoop()

	return a
}

func (a *Agent) writeLoop() {
	for b := range a.writeChan {
		if _, err := a.conn.Write(b); err != nil {
			log.Error("write error(%v)", err)
			break
		}
	}
}

func (a *Agent) Run() {
	for {
		_, err := a.readMsg()
		if err != nil {
			log.Error("read error(%v)", err)
			break
		}
	}
}

func (a *Agent) readMsg() ([]byte, error) {
	return nil, nil
}
