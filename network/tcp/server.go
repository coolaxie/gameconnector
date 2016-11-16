package tcp

import (
	"github.com/coolaxie/gameconnector/log"
	"net"
	"time"
)

type ConnSet map[net.Conn]struct{}

type TCPServer struct {
	Addr  string
	ln    net.Listener
	conns ConnSet

}

func (s *TCPServer) Start() {
	log.Release("tcp server(%v) starting...", s.Addr)
	s.init()
	go s.run()
}

func (s *TCPServer) Close() {
	log.Release("tcp server(%v) closing...", s.Addr)
	s.ln.Close()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
}

func (s *TCPServer) init() {
	s.conns = make(ConnSet)

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatal("listen error(%v)", err)
	}

	s.ln = ln
}

func (s *TCPServer) run() {
	var delay time.Duration
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); !ok || !netErr.Temporary() {
				log.Error("accept error(%v)", err)
				break
			}

			if delay == 0 {
				delay = 5 * time.Millisecond
			} else {
				delay *= 2
				if max := time.Second; delay > max {
					delay = max
				}
			}

			log.Release("accept occur temporary error(%v), delay(%v) retring", err, delay)
			time.Sleep(delay)
		}

		s.conns[conn] = struct {}{}
	}
}
