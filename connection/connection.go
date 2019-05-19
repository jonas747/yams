package connection

import (
	"bytes"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/jonas747/yams/connection/packetmappings"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

type Connection struct {
	manager *ConnectionManager

	id             int
	conn           net.Conn
	byteReaderConn *netByteReader

	ReadError error

	readbuf    []byte
	readCursor int

	writebuf []byte

	compressed bool

	version int
	state   *int32

	uuid uuid.UUID

	writeLock sync.Mutex

	eventHandlerLock sync.RWMutex
	eventHandler     func(packetID packetmappings.YAMPacketID) error
}

func newConn(manager *ConnectionManager, id int, conn net.Conn) *Connection {
	byteReaderConn := &netByteReader{Conn: conn}
	return &Connection{
		id:             id,
		conn:           conn,
		byteReaderConn: byteReaderConn,
		manager:        manager,
		state:          new(int32),
		version:        482,
	}
}

func (c *Connection) Log(f string, args ...interface{}) {
	prefix := "[c:" + strconv.Itoa(c.id) + "] "
	log.Printf(prefix+f, args...)
}

func (c *Connection) GetID() int {
	return c.id
}

func (c *Connection) SetEventHandler(f func(packetID packetmappings.YAMPacketID) error) {
	c.eventHandlerLock.Lock()
	c.eventHandler = f
	c.eventHandlerLock.Unlock()
}

func (c *Connection) Reader() {
	for {
		err := c.readNextPacket()
		if err != nil {
			c.Log("error handling packet: %v", err)
			c.conn.Close()
			break
		}
	}
}

func (c *Connection) SetState(newState State) {
	atomic.StoreInt32(c.state, int32(newState))
	c.Log("set state to %d (%s)", newState, newState)
}

// max 16MB packets
const MaxPacketSize = 0xffffff

var (
	ErrPacketTooBig    = errors.New("Packet is too big (>16MB)")
	ErrPacketSizeSmall = errors.New("Packet size is too small")
)

func (c *Connection) readNextPacket() error {
	c.readCursor = 0

	// read the length of the packet
	rl, err := binary.ReadUvarint(c.byteReaderConn)
	if err != nil {
		c.ReadError = err
		return errors.Wrap(err, "packet_len")
	}
	l := int(rl)

	c.Log("Inc packet: %d", l)

	if l > MaxPacketSize {
		return ErrPacketTooBig
	}

	if l < 1 {
		return ErrPacketSizeSmall
	}

	// read the whole packet into the buffer
	if cap(c.readbuf) < l {
		// need to make a new buffer
		c.readbuf = make([]byte, l)
	} else {
		c.readbuf = c.readbuf[:l]
	}

	_, err = io.ReadFull(c.conn, c.readbuf)
	if err != nil {
		c.ReadError = err
		return err
	}

	return c.handleNextPacketUncompressed()
}

func (c *Connection) handleNextPacketUncompressed() error {
	// read the length of the packet
	packetID, err := c.ReadVarInt()
	if err != nil {
		return err
	}

	return c.emitPacketEvent(packetID)
}

func (c *Connection) emitPacketEvent(packetID int) error {
	state := atomic.LoadInt32(c.state)
	yamID := packetmappings.GetYAMPacketID(c.version, int(state), false, packetID)

	c.eventHandlerLock.RLock()
	cHandler := c.eventHandler
	c.eventHandlerLock.RUnlock()

	c.Log("Handling: 0x%2x (%s), l:%d, data: %#v", packetID, yamID.String(), len(c.readbuf), c.readbuf)

	if cHandler != nil {
		err := cHandler(yamID)
		if err != nil {
			return err
		}
	}

	if targetHandlers, ok := handlers[yamID]; ok {
		for _, v := range targetHandlers {
			err := v(c)
			if err != nil {
				return err
			}
		}
	} else if cHandler == nil {
		c.Log("No handlers for event: 0x%2x (%s), data: %#v", packetID, yamID.String(), c.readbuf)
	}

	return nil
}

func (c *Connection) ReadByte() (byte, error) {
	if c.readCursor >= len(c.readbuf) {
		return 0, io.EOF
	}

	b := c.readbuf[c.readCursor]
	c.readCursor++
	return b, nil
}

func (c *Connection) ReadUInt16() (i uint16, err error) {
	if c.ReadError != nil {
		return 0, c.ReadError
	}

	if c.readCursor+2 > len(c.readbuf) {
		err = errors.New("Out of bounds readuint16")
		c.ReadError = err
		return 0, err
	}

	i = binary.BigEndian.Uint16(c.readbuf[c.readCursor : c.readCursor+2])
	c.readCursor += 2
	return i, nil
}

func (c *Connection) ReadUint64() (i uint64, err error) {
	if c.ReadError != nil {
		return 0, c.ReadError
	}

	if c.readCursor+8 > len(c.readbuf) {
		err = errors.New("Out of bounds readuint16")
		c.ReadError = err
		return 0, err
	}

	i = binary.BigEndian.Uint64(c.readbuf[c.readCursor : c.readCursor+8])
	c.readCursor += 8
	return i, nil
}

func (c *Connection) ReadInt64() (i int64, err error) {
	ui, err := c.ReadUint64()
	return int64(ui), err
}

func (c *Connection) ReadVarInt() (i int, err error) {
	if c.ReadError != nil {
		return 0, c.ReadError
	}

	v, err := binary.ReadUvarint(c)
	if err != nil {
		c.ReadError = err
		return 0, err
	}

	return int(v), nil
}

func (c *Connection) ReadString() (s string, err error) {
	l, err := c.ReadVarInt()
	if err != nil {
		return "", err
	}

	if c.readCursor+l > len(c.readbuf) {
		err = errors.New("Out of bounds string")
		c.ReadError = err
		return "", err
	}

	s = string(c.readbuf[c.readCursor : c.readCursor+l])
	c.readCursor += l
	return s, nil
}

func (c *Connection) WritePacketAsync(id packetmappings.YAMPacketID, components ...PacketField) {
	go func() {
		err := c.WritePacket(id, components...)
		if err != nil {
			c.Log("failed writing packet: %v", err)
		}
	}()
}

func (c *Connection) WritePacket(id packetmappings.YAMPacketID, components ...PacketField) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	buf := bytes.NewBuffer(nil)

	// serialize the packet id
	mcPacketID := packetmappings.GetMCPacketID(c.version, id)
	_, err := WriteVarInt(int(mcPacketID), buf)
	if err != nil {
		return err
	}

	// serialize the components
	for _, comp := range components {
		_, err = comp.MarshalMinecraft(buf)
		if err != nil {
			return err
		}
	}

	// write packet length
	l := buf.Len()
	_, err = WriteVarInt(l, c.conn)
	if err != nil {
		return err
	}

	// write the packet data itself
	_, err = buf.WriteTo(c.conn)
	return err
}
