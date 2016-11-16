package kcp

import "github.com/coolaxie/gameconnector/log"

type KCPServer struct {
	Addr string
}

func (s *KCPServer) Start() {
	log.Release("kcp server(%v) starting...", s.Addr)
}

func (s *KCPServer) Close() {
	log.Release("kcp server(%v) closing...", s.Addr)
}
