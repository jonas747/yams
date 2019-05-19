package game

import (
	"github.com/jonas747/yams/connection"
	"github.com/jonas747/yams/connection/packetmappings"
	"log"
	"sync"
	"sync/atomic"
)

type Game struct {
	mu sync.Mutex

	world   *World
	players []*Player

	entityIDIncr *int32
}

func NewGame() *Game {
	return &Game{
		world:        &World{},
		entityIDIncr: new(int32),
	}
}

func (g *Game) GenEntityID() int32 {
	return atomic.AddInt32(g.entityIDIncr, 1)
}

func (g *Game) AddPlayer(p *Player) {
	g.mu.Lock()
	g.players = append(g.players, p)
	g.mu.Unlock()
}

func (g *Game) Upgrader(c *connection.Connection, username string) {
	p := &Player{
		username: username,
		Conn:     c,
		entityID: int(g.GenEntityID()),
	}

	c.SetEventHandler(p.HandlePacket)

	p.Lock()
	defer p.Unlock()

	// send join game
	err := c.WritePacket(packetmappings.PlayClientJoinGame,
		connection.Int32Field(p.entityID),
		connection.ByteField(1),           // gamemode,  creative
		connection.Int32Field(0),          // dimension, overworld
		connection.ByteField(100),         // max players, unused?
		connection.StringField("default"), // level type (default, flat, largeBiomes, amplified, default_1_1)
		connection.VarIntField(25),        // view distance (2-32)
		connection.BooleanField(false),    // reduced debug info, If true, a Notchian client shows reduced information on the debug screen. For servers in development, this should almost always be false.
	)
	if err != nil {
		log.Printf("error sending join game packet: %v", err)
		return
	}

	// send spawn position
	err = c.WritePacket(packetmappings.PlayClientSpawnPosition, connection.PositionFieldFromComponents(0, 100, 0))
	if err != nil {
		log.Printf("error sending spawn position packet: %v", err)
		return
	}

	// send abilities
	err = c.WritePacket(packetmappings.PlayClientPlayerAbilities,
		connection.ByteField(0x01|0x02|0x04|0x08), // flag
		connection.Float32Field(0.1),              // flying speed
		connection.Float32Field(0.1),              // field of view
	)
	if err != nil {
		log.Printf("error sending abilities packet: %v", err)
		return
	}

	err = c.WritePacket(packetmappings.PlayClientPlayerPositionAndLook,
		connection.Float64Field(0),   // x pos
		connection.Float64Field(100), // y pos
		connection.Float64Field(0),   // z pos

		connection.Float32Field(0), // yaw
		connection.Float32Field(0), // pitch

		connection.ByteField(0),      // flags
		connection.VarIntField(1337), // teleport confirm
	)

	if err != nil {
		log.Printf("error sending abilities packet: %v", err)
		return
	}

	g.AddPlayer(p)
}
