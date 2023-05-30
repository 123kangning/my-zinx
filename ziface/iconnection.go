package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	// Start 启动连接，让当前的连接准备开始工作
	Start()
	// Stop 停止连接 结束当前连接的工作
	Stop()
	// GetTCPConnection 获取当前连接所绑定的socket connection
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接的连接ID
	GetConnID() uint32
	// RemoteAddr 获取远程客户端的 TCP状态（ip port）
	RemoteAddr() net.Addr
	// SendMsg 发送数据，将数据发送给远程客户端，先封包，再发送
	SendMsg(id uint32, data []byte) error
	// SetProperty 设置链接属性
	SetProperty(key string, value interface{})
	// GetProperty 获取链接属性
	GetProperty(key string) (interface{}, error)
	// RemoveProperty 移除链接属性
	RemoveProperty(key string)
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
