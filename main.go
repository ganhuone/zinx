package main

import (
	"fmt"
	"time"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (h *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	fmt.Println("recv from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(request.GetMsgId(), []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnBegin(conn ziface.IConnection) {
	fmt.Println("--> DoConnBegin")
	conn.SetProperty("OpenTime", time.Now().Format("2006-01-02 15:04:05"))
	err := conn.SendMsg(202, []byte("Do connection BEGIN"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnLost(conn ziface.IConnection) {
	fmt.Println("--> DoConnLost")
	openTime, err := conn.GetProperty("OpenTime")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("open time is ", openTime)
	}

	fmt.Println("conn id = ", conn.GetConnID(), " is Lost...")
}

func main() {
	s, err := znet.NewServer("[zinx V0.1]")
	if err != nil {
		fmt.Println(err)
		return
	}

	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
