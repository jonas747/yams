package connection

import (
	"encoding/binary"
	"io"
)

type PacketComponent interface {
	MarshalMinecraft(w io.Writer) (n int, err error)
}

func WriteVarInt(i int, w io.Writer) (n int, err error) {
	buf := make([]byte, 5)
	n = binary.PutUvarint(buf, uint64(uint32(i)))
	return w.Write(buf[:n])
}

type VarIntComponent int

func (i VarIntComponent) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteVarInt(int(i), w)
}

func WriteString(s string, w io.Writer) (n int, err error) {
	n, err = WriteVarInt(len(s), w)
	if err != nil {
		return n, err
	}

	n2, err := w.Write([]byte(s))
	return n + n2, err
}

type StringComponent string

func (s StringComponent) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteString(string(s), w)
}
