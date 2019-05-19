package connection

import (
	"io"
	"net"
)

type netByteReader struct {
	net.Conn
	buf []byte
}

func (n *netByteReader) ReadByte() (b byte, err error) {
	if n.buf == nil {
		n.buf = make([]byte, 1)
	}

	_, err = io.ReadFull(n.Conn, n.buf)
	return n.buf[0], err
}
