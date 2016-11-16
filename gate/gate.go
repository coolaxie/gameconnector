package gate

import "github.com/coolaxie/gameconnector/network/tcp"

type Gate struct {
	TCPAddrs []string
}

func (g *Gate) Run(closeSig chan bool) {
	tcpServers := make([]*tcp.TCPServer, len(g.TCPAddrs))
	for i := 0; i < len(g.TCPAddrs); i++ {
		tcpServers[i] = new(tcp.TCPServer)
		tcpServers[i].Addr = g.TCPAddrs[i]
		tcpServers[i].Start()
	}

	<-closeSig

	for i := 0; i < len(tcpServers); i++ {
		tcpServers[i].Close()
	}
}
