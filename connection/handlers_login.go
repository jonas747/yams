package connection

import (
	"github.com/google/uuid"
	"github.com/jonas747/yams/connection/packetmappings"
)

func init() {
	RegisterHandler(packetmappings.LoginServerLoginStart, handleLoginStart)
}

func handleLoginStart(c *Connection) error {
	username, err := c.ReadString()
	if err != nil {
		return err
	}

	c.uuid, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	c.SetState(StatePlay)

	err = c.WritePacket(packetmappings.LoginClientLoginSuccess, StringField(c.uuid.String()), StringField(username))
	if err != nil {
		return err
	}

	go c.manager.Upgrader(c, username)

	return nil
}
