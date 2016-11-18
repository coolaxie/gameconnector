package tcp

import (
	"encoding/binary"
	"errors"
	"github.com/coolaxie/gameconnector/log"
	"io"
	"net"
	"sync"
)

type TCPConn struct {
	sync.Mutex
	conn      net.Conn
	writeChan chan []byte
	closeFlag bool

	lenMsgLen    int
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

func NewTCPConn(conn net.Conn) *TCPConn {
	tcpConn := new(TCPConn)
	tcpConn.conn = conn
	tcpConn.writeChan = make(chan []byte, 2000)

	tcpConn.lenMsgLen = 4
	tcpConn.minMsgLen = 1
	tcpConn.maxMsgLen = 4096
	tcpConn.littleEndian = false

	go func() {
		for b := range tcpConn.writeChan {
			if b == nil {
				break
			}

			if _, err := conn.Write(b); err != nil {
				break
			}
		}

		conn.Close()
		tcpConn.Lock()
		tcpConn.closeFlag = true
		tcpConn.Unlock()
	}()

	return tcpConn
}

func (tcpConn *TCPConn) Destroy() {
	tcpConn.Lock()
	defer tcpConn.Unlock()

	tcpConn.doDestroy()
}

func (tcpConn *TCPConn) Close() {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag {
		return
	}

	tcpConn.Write(nil)
	tcpConn.closeFlag = true
}

func (tcpConn *TCPConn) ReadMsg() ([]byte, error) {
	var b [4]byte //max len for msg len is 4
	bufMsgLen := b[:tcpConn.lenMsgLen]

	if _, err := io.ReadFull(tcpConn.conn, bufMsgLen); err != nil {
		return nil, err
	}

	var msgLen uint32
	switch tcpConn.lenMsgLen {
	case 1:
		msgLen = uint32(bufMsgLen[0])
	case 2:
		if tcpConn.littleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if tcpConn.littleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}

	if msgLen > tcpConn.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < tcpConn.minMsgLen {
		return nil, errors.New("message too short")
	}

	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(tcpConn.conn, msgData); err != nil {
		return nil, err
	}

	return msgData, nil
}

func (tcpConn *TCPConn) WriteMsg(data []byte) error {
	msgLen := uint32(len(data))
	if msgLen > tcpConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < tcpConn.minMsgLen {
		return errors.New("message too short")
	}

	msg := make([]byte, uint32(tcpConn.lenMsgLen)+msgLen)
	switch tcpConn.lenMsgLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if tcpConn.littleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if tcpConn.littleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}

	copy(msg[tcpConn.lenMsgLen:], data)
	tcpConn.Write(msg)

	return nil
}

func (tcpConn *TCPConn) Write(data []byte) {
	tcpConn.Lock()
	defer tcpConn.Unlock()
	if tcpConn.closeFlag || data == nil {
		return
	}

	tcpConn.doWrite(data)
}

func (tcpConn *TCPConn) RemoteAddr() string {
	return tcpConn.conn.RemoteAddr().String()
}

func (tcpConn *TCPConn) doWrite(data []byte) {
	if len(tcpConn.writeChan) == cap(tcpConn.writeChan) {
		log.Debug("close conn because channel full")
		tcpConn.doDestroy()
		return
	}

	tcpConn.writeChan <- data
}

func (tcpConn *TCPConn) doDestroy() {
	tcpConn.conn.(*net.TCPConn).SetLinger(0)
	tcpConn.conn.Close()

	if !tcpConn.closeFlag {
		close(tcpConn.writeChan)
		tcpConn.closeFlag = true
	}
}
