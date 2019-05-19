package connection

import (
	"github.com/jonas747/yams/connection/packetmappings"
)

type HandlerFunc func(c *Connection) error

var handlers = make(map[packetmappings.YAMPacketID][]HandlerFunc)

func init() {
	RegisterHandler(packetmappings.HandshakingServerHandshake, handleHandshake)
}

func RegisterHandler(packetID packetmappings.YAMPacketID, handler HandlerFunc) {
	handlers[packetID] = append(handlers[packetID], handler)
}

func handleHandshake(c *Connection) error {
	protcocolVersion, _ := c.ReadVarInt()
	serverAddr, _ := c.ReadString()
	serverPort, _ := c.ReadUInt16()
	nextState, err := c.ReadVarInt()

	c.Log("Got handshake! v:%d, addr:%s, port:%d, nextState:%d, err:%v", protcocolVersion, serverAddr, serverPort, nextState, err)

	if err != nil {
		return err
	}

	c.SetState(State(nextState))

	return nil
}
