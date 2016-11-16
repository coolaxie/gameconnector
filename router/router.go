package router

import (
	"github.com/coolaxie/gameconnector/log"
	"github.com/coolaxie/gameconnector/network/tcp"
)

type Router struct {
	TCPAddr string
}

func (r *Router) Run() {
	if r.TCPAddr == "" {
		log.Fatal("must give router tcp addr")
	}

	tcpServer := new(tcp.TCPServer)
	tcpServer.Addr = r.TCPAddr
	tcpServer.Start()

	

}