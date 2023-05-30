package main

import (
	"fmt"
	"my-zinx/utils"
	"my-zinx/znet"
	"net"
	"time"
)

/*
模拟客户端
*/
func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	utils.Init()
	for {
		//发封包message消息
		dp := znet.NewDataPack()
		msg1, _ := dp.Pack(znet.NewMessage(0, []byte("Zinx V0.8 Client0 Test Message")))
		_, err := conn.Write(msg1)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//将headData字节流 拆包到msg中
		msg, err := dp.UnPack(conn)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		fmt.Println("receive server message , ID = ", msg.GetMsgId(), " DataLen = ", msg.GetMsgLen(), " Data = ", string(msg.GetData()))

		time.Sleep(1 * time.Second)
	}
}
