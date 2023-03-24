package znet

import (
	"fmt"
	"my-zinx/ziface"
)

/*
MsgHandle 消息处理模块的实现
*/
type MsgHandle struct {
	//存放每一个MsgID所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandle 创建MsgHandle对象
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	//1.通过request中找到 msgHandler
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Printf("api msgID = %d is not found ! Need register!\n", request.GetMsgId())
		return
	}
	//2.根据msgHandler调度对应router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1.判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic(fmt.Sprintf("repeat api , msgID = %d ", msgID))
	}
	//2.添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %d success \n", msgID)
}
