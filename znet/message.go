package znet

type Message struct {
	Id      uint32 //消息的ID
	DataLen uint32 //消息的长度
	Data    []byte //消息的内容
}

// NewMessage 创建一个Message消息包
func NewMessage(id uint32, data []byte) *Message {
	return &Message{DataLen: uint32(len(data)), Id: id, Data: data}
}

// GetMsgId 获取消息ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// GetMsgLen 获取消息长度
func (msg *Message) GetMsgLen() uint32 {
	return msg.DataLen
}

// GetData 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

// SetMsgId 设置消息ID
func (msg *Message) SetMsgId(id uint32) {
	msg.Id = id
}

// SetMsgLen 设置消息长度
func (msg *Message) SetMsgLen(msgLen uint32) {
	msg.DataLen = msgLen
}

// SetData 设置消息内容
func (msg *Message) SetData(data []byte) {
	msg.Data = data
}
