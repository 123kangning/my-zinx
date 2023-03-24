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

// PingRouter test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// PreHandle Test
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	conn := request.GetConnection().GetTCPConnection()
	if _, err := conn.Write([]byte("before ping..\n")); err != nil {
		fmt.Println("Call Back Before Ping error ", err)
	}
}

// Handle 处理conn主业务的方法
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	conn := request.GetConnection().GetTCPConnection()
	if _, err := conn.Write(append(request.GetData(), []byte("ping..\n")...)); err != nil {
		fmt.Println("Call Back Ping error ", err)
	}
}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	conn := request.GetConnection().GetTCPConnection()
	if _, err := conn.Write([]byte("after ping..\n")); err != nil {
		fmt.Println("Call Back after Ping error ", err)
	}
}

func main() {
	utils.Init()
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer()
	//给当前框架添加一个自定义的router(暂时只能注册一个路由)
	s.AddRouter(0, &PingRouter{})
	//启动server
	s.Serve()
}
