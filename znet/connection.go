package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/ziface"
)

type Connection struct {
	TcpServer ziface.IServer

	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	msgChan chan []byte

	ExitChan chan bool

	MsgHandler ziface.IMsgHandler

	property map[string]interface{}

	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connId uint32, msgHandler ziface.IMsgHandler) *Connection {

	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connId,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnManager().Add(c)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())

	defer c.Stop()

	for {

		dp := NewDataPack()

		headData := make([]byte, dp.GetHeadLen())

		_, err := io.ReadFull(c.CetTcpConnection(), headData)
		if err != nil {
			fmt.Println(err)
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println(err)
			break
		}

		if msg.GetDataLen() > 0 {
			data := make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(c.CetTcpConnection(), data)
			if err != nil {
				fmt.Println(err)
				break
			}

			msg.SetData(data)
		}

		c.MsgHandler.SendMsgToTaskQueue(&Request{
			conn: c,
			msg:  msg,
		})
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Writer is exit, remote addr is", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			_, err := c.CetTcpConnection().Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()...ConnID = ", c.ConnID)

	go c.StartReader()

	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop()...ConnID = ", c.ConnID)

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.TcpServer.CallOnConnStop(c)

	c.Conn.Close()

	c.ExitChan <- true

	c.TcpServer.GetConnManager().Remove(c)

	close(c.msgChan)
	close(c.ExitChan)

}

func (c *Connection) CetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		return err
	}

	c.msgChan <- binaryMsg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	value, ok := c.property[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
