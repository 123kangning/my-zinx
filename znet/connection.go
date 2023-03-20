package znet

import (
	"fmt"
	"my-zinx/ziface"
	"net"
)

/*
连接模块
*/
type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态
	IsClosed bool
	//当前连接所绑定的处理业务方法API
	handleAPI ziface.HandleFunc
	//告知当前连接已经退出/停止的 channel
	ExitChan chan bool
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		IsClosed:  false,
		handleAPI: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
	return c
}

// StartReader 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//读取客户端的数据到buf中，最大 512 字节的存储
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}
		//调用当前连接所绑定的HandleApi
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID ", c.ConnID, " handler is error", err)
			continue
		}
	}
}

// Start 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接的写数据的业务

}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)
	//如果当前连接已经关闭
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	//关闭socket连接
	c.Conn.Close()
	//回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前连接所绑定的socket connection
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接的连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的 TCP状态（ip port）
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// Send 发送数据，将数据发送给远程客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
