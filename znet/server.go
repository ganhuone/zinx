package znet

import (
	"fmt"
	"net"

	"github.com/ganhuone/zinx/utils"
	"github.com/ganhuone/zinx/ziface"
)

type Server struct {
	Name string

	IPVersion string

	IP string

	Port int

	MsgHandler ziface.IMsgHandler

	ConnMagager ziface.IConnManager

	OnConnStart func(conn ziface.IConnection)

	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {

	fmt.Printf("[Start] Server - %s Listenner at IP: %s, Port %d, is starting\n", s.Name, s.IP, s.Port)

	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("resolve tcp addr error: ", err)
		return
	}

	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		return
	}

	fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning...")

	connId := 1

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err: ", err)
			continue
		}

		if s.ConnMagager.Len() >= utils.GlobalObject.MaxConn {
			fmt.Println("Too many connection MaxConn = ", utils.GlobalObject.MaxConn)
			conn.Close()
			continue
		}

		dealConn := NewConnection(s, conn, uint32(connId), s.MsgHandler)

		dealConn.Start()

		connId++

	}

}

func (s *Server) Stop() {
	fmt.Println("[STOP] server - ", s.Name)
	s.ConnMagager.ClearConn()

	s.MsgHandler.StopWorkerPool()
}

func (s *Server) Serve() {
	s.MsgHandler.StartWorkerPool()

	go s.Start()

	select {}

}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func NewServer(name string) (ziface.IServer, error) {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnMagager: NewConnManager(),
	}

	return s, nil
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnMagager
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--> call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--> call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
