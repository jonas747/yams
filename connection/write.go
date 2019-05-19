package connection

import (
	"encoding/binary"
	"io"
)

type PacketField interface {
	MarshalMinecraft(w io.Writer) (n int, err error)
}

func WriteVarInt(i int, w io.Writer) (n int, err error) {
	buf := make([]byte, 5)
	n = binary.PutUvarint(buf, uint64(uint32(i)))
	return w.Write(buf[:n])
}

type VarIntField int

func (i VarIntField) MarshalMinecraft(w io.Writer) (n int, err error) {
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

type StringField string

func (s StringField) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteString(string(s), w)
}

func WriteInt64(i int64, w io.Writer) (n int, err error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return w.Write(buf)
}

type Int64Field int64

func (i Int64Field) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteInt64(int64(i), w)
}
