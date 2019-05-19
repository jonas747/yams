package connection

import (
	"encoding/binary"
	"io"
	"math"
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

type ByteField byte

func (i ByteField) MarshalMinecraft(w io.Writer) (n int, err error) {
	return w.Write([]byte{byte(i)})
}

func WriteInt32(i int32, w io.Writer) (n int, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return w.Write(buf)
}

type Int32Field int32

func (i Int32Field) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteInt32(int32(i), w)
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

type BooleanField bool

func (b BooleanField) MarshalMinecraft(w io.Writer) (n int, err error) {
	if b {
		return w.Write([]byte{1})
	}
	return w.Write([]byte{0})
}

type PositionField int64

func PositionFieldFromComponents(x, y, z int) PositionField {
	return PositionField(((int64(x) & 0x3FFFFFF) << 38) | ((int64(z) & 0x3FFFFFF) << 12) | (int64(y) & 0xFFF))
}

func (p PositionField) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteInt64(int64(p), w)
}

type Float32Field float32

func WriteFloat32(f float32, w io.Writer) (n int, err error) {
	bits := math.Float32bits(f)
	return WriteInt32(int32(bits), w)
}

func (f Float32Field) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteFloat32(float32(f), w)
}

type Float64Field float64

func WriteFloat64(f float64, w io.Writer) (n int, err error) {
	bits := math.Float64bits(f)
	return WriteInt64(int64(bits), w)
}

func (f Float64Field) MarshalMinecraft(w io.Writer) (n int, err error) {
	return WriteFloat64(float64(f), w)
}
