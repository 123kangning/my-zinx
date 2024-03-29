package znet

import "my-zinx/ziface"

// BaseRouter 实现Router时，先嵌入BaseRouter这个基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

// 这里之所以BaseRouter的方法都为空，是因为有的Router不希望有PreHandle和PostHandle这两个业务，
//所以Router全部继承BaseRouter的好处就是，不需要实现PreHandle和PostHandle

// PreHandle 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// Handle 处理conn主业务的方法
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
