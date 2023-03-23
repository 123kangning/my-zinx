package znet

import (
	"fmt"
	"my-zinx/utils"
	"net"
	"testing"
)

func TestDataPack_Pack(t *testing.T) {
	//buf := []byte("this is my test")
	//msg := Message{Id: 1, DataLen: uint32(len(buf)), Data: buf}
	//ans, err := DataPack.Pack(msg)
	/*
		模拟服务器
	*/
	utils.Init()
	//1.创建socketTCP
	listenner, err := net.Listen("tcp", ":7777")
	//2.从客户端读取数据，拆包处理
	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error , ", err)
			}
			go func(conn net.Conn) {
				dp := DataPack{}
				for {

					msg, err := dp.UnPack(conn)
					if err != nil {
						fmt.Println("server unPack error ", err)
						return
					}

					//本次读取结束
					//fmt.Printf("MaxPackageSize = %d, msg.Id = %d , msg.DataLen = %d,msg = %s\n", utils.GlobalObject.MaxPackageSize, msg.GetMsgId(), msg.GetMsgLen(), msg.GetData())
					fmt.Println("msg.Id = ", msg.GetMsgId(), "msg.DataLen = ", msg.GetMsgLen(), " msg.data = ", msg.GetData())
				}

			}(conn)
		}
	}()
	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", ":7777")
	if err != nil {
		fmt.Println("client dial err ", err)
		return
	}
	//创建一个封包对象
	dp := NewDataPack()
	for i := 0; i < 22700; i++ {
		//封装第一个msg包
		msg1 := &Message{
			Id:      1,
			DataLen: 9,
			Data:    []byte("111111111"),
		}
		//封装第二个msg包
		msg2 := &Message{
			Id:      2,
			DataLen: 9,
			Data:    []byte("000000000"),
		}
		//将两个包黏在一起
		data1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 error ", err)
			return
		}
		data2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client pack msg2 error ", err)
			return
		}
		//一次性发送给服务端
		_, err = conn.Write(append(data1, data2...))
		if err != nil {
			fmt.Println("con write error ", err)
			return
		}
	}
	//发送之后客户端阻塞
	select {}
}
