package game

import (
	"github.com/jonas747/yams/connection"
	"github.com/jonas747/yams/connection/packetmappings"
	"sync"
)

type Player struct {
	sync.Mutex

	Conn *connection.Connection

	entityID int
	username string
}

func (p *Player) Log(f string, args ...interface{}) {
	prefix := "[p:" + p.username + "] "
	p.Conn.Log(prefix+f, args...)
}

func (p *Player) GetUsername() string {
	return p.username
}

func (p *Player) GetEntityID() int {
	return p.entityID
}

func (p *Player) HandlePacket(packetID packetmappings.YAMPacketID) error {
	p.Lock()
	defer p.Unlock()

	switch packetID {
	default:
		p.Log("No handler for packet: %02X (%s)", packetID, packetID.String())
	}

	return nil
}

func (p *Player) Tick(t int) {

}
