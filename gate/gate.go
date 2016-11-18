package gate

import "github.com/coolaxie/gameconnector/network/tcp"

type Gate struct {
	TCPAddr string
	KCPAddr string
}

func (g *Gate) Run(closeSig chan bool) {
	var tcpServer *tcp.TCPServer
	if g.TCPAddr != "" {
		tcpServer = new(tcp.TCPServer)
		tcpServer.Addr = g.TCPAddr
		tcpServer.NewAgent = newTCPAgent
	}

	if tcpServer != nil {
		tcpServer.Start()
	}

	<-closeSig

	if tcpServer != nil {
		tcpServer.Close()
	}
}


