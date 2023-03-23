package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"my-zinx/utils"
	"my-zinx/ziface"
	"net"
)

// DataPack 封包、拆包的具体模块
type DataPack struct{}

// NewDataPack 封包拆包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// GetHeadLen 获取包的头的长度方法？？？
func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节)+ID uint32(4字节)
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuf := bytes.Buffer{}
	//将dataLen写入dataBuf中
	if err := binary.Write(&dataBuf, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将msgID写入dataBuf中
	if err := binary.Write(&dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data数据写入dataBuf中
	if err := binary.Write(&dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

// UnPack 拆包方法，传入 net.Conn 对象，返回读取到的Message对象
func (dp *DataPack) UnPack(c net.Conn) (ziface.IMessage, error) {
	//读取客户端的MsgHead二进制流 8个字节
	headData := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(c, headData); err != nil {
		fmt.Println("read msg head error ", err)
		return nil, err
	}
	//创建一个输入的二进制数据的ioReader
	dataBuf := bytes.NewReader(headData)
	msg := &Message{}
	//读取dataLen
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读取MsgID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && utils.GlobalObject.MaxPackageSize < msg.DataLen {
		return nil, errors.New(fmt.Sprintf("msg too large to read %d %d", utils.GlobalObject.MaxPackageSize, msg.DataLen))
	}
	//根据dataLen 再次读取Data ，放在msg.Data中
	if msg.GetMsgLen() > 0 {
		data := make([]byte, msg.GetMsgLen())
		_, err := io.ReadFull(c, data)
		if err != nil {
			fmt.Println("server unPack data error ", err)
			return nil, errors.New("server unPack data error ")
		}
		msg.SetData(data)
	}
	return msg, nil
}
