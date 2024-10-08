package ziface

type IServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint32, router IRouter)

	GetConnManager() IConnManager

	SetOnConnStart(func(conn IConnection))

	SetOnConnStop(func(conn IConnection))

	CallOnConnStart(conn IConnection)

	CallOnConnStop(conn IConnection)
}
