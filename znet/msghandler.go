package znet

import (
	"fmt"
	"math/rand"

	"github.com/ganhuone/zinx/utils"
	"github.com/ganhuone/zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter

	TaskQueue []chan ziface.IRequest

	WorkerPoolSize uint32

	isStop chan bool
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		isStop:         make(chan bool, 1),
	}
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := m.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), "is NOT FOUND")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		fmt.Println("is register")
		return
	}

	m.Apis[msgId] = router
	fmt.Println("success")

}

func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		workerTaskQueue := make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		m.TaskQueue[i] = workerTaskQueue

		go m.StartOneWorker(i, workerTaskQueue)
	}
}

func (m *MsgHandler) StopWorkerPool() {
	for _, v := range m.TaskQueue {
		close(v)
	}

	m.isStop <- true
}

func (m *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker id = ", workerID, " is start")
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		case <-m.isStop:
			return
		}
	}
}

func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerId := rand.Intn(int(m.WorkerPoolSize) - 1)

	m.TaskQueue[workerId] <- request

	fmt.Println("Add Conn = ", request.GetConnection().GetConnID(), " request MsgId = ", request.GetMsgId(), " to workerId = ", workerId)
}
