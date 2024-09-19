package znet

import (
	"errors"
	"fmt"
	"sync"
	"github.com/ganhuone/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn

	fmt.Println("connId = ", conn.GetConnID(), " add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())

	fmt.Println("connId = ", conn.GetConnID(), " remove to ConnManager successfully: conn num = ", c.Len())

}

func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	conn, ok := c.connections[connId]
	if ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()

		delete(c.connections, connID)
	}

	fmt.Println("clear all connection success, conn num = ", c.Len())
}
