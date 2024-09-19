package utils

import (
	"encoding/json"
	"os"
	"github.com/ganhuone/zinx/ziface"
)

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	GlobalObject.Reload()
}

type GlobalObj struct {
	TcpServcer ziface.IServer //当前全局server对象

	Host string //监听的ip

	TcpPort int //监听的端口

	Name string //服务名称

	Version string //服务版本号

	MaxConn int //允许的最大链接数

	MaxPackageSize uint32 //数据包最大值

	WorkerPoolSize uint32 //当前业务工作worker池的Goroutine数量(干活的人数)

	MaxWorkerTaskLen uint32 //单个worker任务队列(有多少活要干)
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}

}
