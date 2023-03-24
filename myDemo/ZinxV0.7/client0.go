package main

import (
	"fmt"
	"my-zinx/utils"
	"my-zinx/znet"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1000000)
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("conn error ", err)
		return
	}
	utils.Init()
	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("ZinxV0.6 client0 test Message")))
		if err != nil {
			fmt.Println("Pack error ", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}
		//server应该给我们恢复一个Message
		msg, err := dp.UnPack(conn)
		fmt.Println("receive server message , ID = ", msg.GetMsgId(), " DataLen = ", msg.GetMsgLen(), " Data = ", string(msg.GetData()))
		time.Sleep(time.Second)
	}

}
