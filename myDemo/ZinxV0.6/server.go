package main

import (
	"fmt"
	"my-zinx/utils"
	"my-zinx/ziface"
	"my-zinx/znet"
)

/*
	基于 Zinx 框架来开发的服务端应用程序
*/

// HelloRouter test自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// Handle 处理conn主业务的方法
func (pr *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Printf("msg.Id = %d , msg.Data = %s\n", request.GetMsgId(), request.GetData())
	err := request.GetConnection().SendMsg(10, []byte("HelloRouter Zinx\n"))
	if err != nil {
		fmt.Println("Call Back Ping error ", err)
	}
}

// PingRouter test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Handle 处理conn主业务的方法
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	fmt.Printf("msg.Id = %d , msg.Data = %s\n", request.GetMsgId(), request.GetData())
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...\n"))
	if err != nil {
		fmt.Println("Call Back Ping error ", err)
	}
}

func main() {
	utils.Init()
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer()
	//给当前框架添加一个自定义的router(暂时只能注册一个路由)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	//启动server
	s.Serve()
}
