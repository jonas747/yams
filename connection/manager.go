package connection

import (
	"log"
	"net"
	"sync"
)

type ConnectionManager struct {
	activeConnection []*Connection
	connectionsMU    sync.Mutex

	// called after a sucessfull login
	Upgrader func(c *Connection, username string)

	connCounter int
}

func (cm *ConnectionManager) Listen(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Server launched on address", addr)

	// go c.keepAlive()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Error accepting connection:", err)
		} else {
			cm.connCounter += 1
			go cm.handleConnection(conn, cm.connCounter)
		}
	}
}

func (cm *ConnectionManager) handleConnection(conn net.Conn, id int) {
	yamConn := newConn(cm, id, conn)
	yamConn.Log("connection opened")

	cm.connectionsMU.Lock()
	cm.activeConnection = append(cm.activeConnection, yamConn)
	cm.connectionsMU.Unlock()

	yamConn.Reader()

	yamConn.Log("connection closed")

	cm.connectionsMU.Lock()
	for i, v := range cm.activeConnection {
		if v.GetID() == id {
			cm.activeConnection = append(cm.activeConnection[:i], cm.activeConnection[i+1:]...)
		}
	}
	cm.connectionsMU.Unlock()
}
