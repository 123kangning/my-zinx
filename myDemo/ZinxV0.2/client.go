package main

import (
	"fmt"
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
	for {
		cnt, err := conn.Write([]byte("hello zinxV0.2"))
		if err != nil {
			fmt.Println("write conn error ", err)
			return
		} else {
			fmt.Println("write bytes ", cnt)
		}
		buf := make([]byte, 512)
		cnt, err = conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
		}
		fmt.Println("read buf = ", string(buf[:cnt]))
		time.Sleep(time.Second)
	}

}
