package znet

import (
	"errors"
	"fmt"
	"my-zinx/ziface"
	"net"
)

/*
Connection 连接模块
*/
type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态
	IsClosed bool
	//告知当前连接已经退出/停止的 channel
	ExitChan chan bool
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API的关系
	MsgHandler ziface.IMsgHandle
}

// NewConnection 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: NewMsgHandle(),
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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		//创建一个拆包解包对象
		dp := NewDataPack()
		//拆包，得到msg.ID 和 msgDataLen放在msg消息中
		msg, err := dp.UnPack(c.Conn)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		//得到当前conn连接的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		go c.MsgHandler.DoMsgHandler(&req)

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

// SendMsg 发送数据，将数据发送给远程客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection close ")
	}
	//将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("pack error ", err, " msgId = ", msgId)
		return errors.New("pack error ")
	}
	//将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("write msg error ", err, " msgId = ", msgId)
		return errors.New("write msg error ")
	}
	return nil
}
