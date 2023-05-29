package ziface

/*
IMsgHandle 消息管理抽象层
*/
type IMsgHandle interface {
	// DoMsgHandler 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)
	// StartWorkerPool 启动worker工作池
	StartWorkerPool()
	// SendMsgToTaskQueue 将请求交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
