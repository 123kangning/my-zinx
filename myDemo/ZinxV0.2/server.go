package main

import "my-zinx/znet"

/*
	基于 Zinx 框架来开发的服务端应用程序
*/

func main() {
	//创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx V0.2]")
	//启动server
	s.Serve()
}
