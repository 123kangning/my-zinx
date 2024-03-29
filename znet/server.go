package znet

import (
	"fmt"
	"my-zinx/utils"
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
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API的关系,一般是用户绑定的
	MsgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

// Start 启动服务器
func (s *Server) Start() {
	//以后可以统一打到日志文件中
	fmt.Printf("[Zinx] Server Name :%s Server Listenner ar IP:%s,Port %d,is atarting\n", s.Name, s.IP, s.Port)
	go func() {
		go s.MsgHandler.StartWorkerPool()
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
			//设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			//启动当前的连接业务处理
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息，进行停止或回收
	s.ConnMgr.ClearConn()
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
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		TCPVersion: "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success")
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
