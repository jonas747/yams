package connection

import (
	"encoding/json"
	"github.com/jonas747/yams/connection/packetmappings"
)

func init() {
	RegisterHandler(packetmappings.StatusServerRequest, handleStatusRequest)
	RegisterHandler(packetmappings.StatusServerPing, handleStatusPing)
}

/*

{
    "version": {
        "name": "1.8.7",
        "protocol": 47
    },
    "players": {
        "max": 100,
        "online": 5,
        "sample": [
            {
                "name": "thinkofdeath",
                "id": "4566e69f-c907-48ee-8d71-d7ba5aa00d20"
            }
        ]
    },
    "description": {
        "text": "Hello world"
    },
    "favicon": "data:image/png;base64,<data>"
}

*/

type StatusResponse struct {
	Version     *StatusResponseVersion     `json:"version"`
	Players     *StatusResponsePlayers     `json:"players"`
	Description *StatusResponseDescription `json:"description"`
	Favicon     *string                    `json:"favicon,omitempty"`
}

type StatusResponseVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type StatusResponsePlayers struct {
	Max    int                            `json:"max"`
	Online int                            `json:"online"`
	Sample []*StatusResponsePlayersSample `json:"sample"`
}

type StatusResponsePlayersSample struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type StatusResponseDescription struct {
	Text string `json:"text"`
}

func handleStatusRequest(c *Connection) error {
	c.Log("Got status request!")

	resp := &StatusResponse{
		Version: &StatusResponseVersion{
			Name:     "SUPER FANCY VERSION",
			Protocol: 482,
		},
		Players: &StatusResponsePlayers{
			Max:    0xffff,
			Online: 0xffff - 1,
		},
		Description: &StatusResponseDescription{
			Text: "This is a super fancy server that dosen't work at all! it's written in go tough...",
		},
	}

	encodedResp, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return c.WritePacket(packetmappings.StatusClientResponse, StringField(encodedResp))
}

func handleStatusPing(c *Connection) error {
	payload, err := c.ReadInt64()
	if err != nil {
		return err
	}

	c.Log("Got status ping! %d", payload)
	return c.WritePacket(packetmappings.StatusClientPong, Int64Field(payload))
}
