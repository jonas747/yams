package connection

import (
	"log"
	"net"
	"sync"
)

type ConnectionManager struct {
	activeConnection []*Connection
	connectionsMU    sync.Mutex

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
	log.Println("New connection ", id)

	yamConn := newConn(cm, id, conn)
	cm.connectionsMU.Lock()
	cm.activeConnection = append(cm.activeConnection, yamConn)
	cm.connectionsMU.Unlock()

	yamConn.Reader()

	log.Println("Connection closed ", id)
}
