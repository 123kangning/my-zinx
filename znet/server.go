package znet

import (
	"errors"
	"fmt"
	"my-zinx/ziface"
	"net"
)

// Server IServer的接口实现
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	TCPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

// CallBackToClient 定义当前连接所绑定的 handleApi(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn handle] CallBackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner ar IP:%s,Port %d,is atarting\n", s.IP, s.Port)
	go func() {
		//1 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.TCPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.TCPVersion, addr)
		if err != nil {
			fmt.Println("Listener ", s.TCPVersion, " error: ", err)
			return
		}
		fmt.Println("start zinx success ", s.Name, " Listening...")

		var cid uint32 = 0
		//3 阻塞的等待客户端连接，处理客户端连接业务（读写）
		for {
			//如果有客户端连接，阻塞返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息，进行停止或回收
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server而的服务功能
	s.Start()
	//TODO 做一些启动服务器之后的额外业务，将来以此丰富扩展框架

	//阻塞状态
	select {}
}

// NewServer 初始化Server模块方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       name,
		TCPVersion: "tcp4",
		IP:         "0.0.0.0",
		Port:       8999,
	}
	return s
}
