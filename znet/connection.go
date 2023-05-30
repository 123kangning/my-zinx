package znet

import (
	"errors"
	"fmt"
	"my-zinx/utils"
	"my-zinx/ziface"
	"net"
	"sync"
)

/*
Connection 连接模块
*/
type Connection struct {
	//当前conn属于哪个server，在conn初始化的时候添加即可
	TcpServer ziface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//连接的ID
	ConnID uint32
	//当前的连接状态
	IsClosed bool
	//告知当前连接已经退出/停止的 channel(由reader告知writer退出)
	ExitChan chan bool
	//用于读，写Goroutine之间的消息通信的管道 缓冲区大小为10
	msgChan chan []byte
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API的关系，一般是从serve中直接拿过来
	MsgHandler ziface.IMsgHandle
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		IsClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte, 10),
		MsgHandler: handler,
		property:   make(map[string]interface{}),
	}
	//将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)
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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// StartWriter 连接的写业务处理,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " Writer is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Write也要退出
			return
		}
	}
}

// Start 启动连接，让当前的连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)
	//启动从当前连接的读数据的业务
	go c.StartReader()
	//启动从当前连接的写数据的业务
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止连接 结束当前连接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID)

	//如果当前连接已经关闭
	if c.IsClosed {
		return
	}
	c.IsClosed = true
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)
	//将链接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)
	//关闭socket连接
	c.Conn.Close()
	//告知Writer关闭(close(c.ExitChan)之后，c.ExitChan可读，就不用这一步，先写上)
	c.ExitChan <- true
	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
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
	//将数据发送给Writer Goroutine
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
