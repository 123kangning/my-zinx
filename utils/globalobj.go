package utils

import (
	"encoding/json"
	"my-zinx/ziface"
	"os"
)

/*
GlobalObj
存储一切有关Zinx框架的全局参数，供其他模块使用
一切参数可以通过zinx.json由用户进行配置
*/
type GlobalObj struct {
	//Server
	TcpServer ziface.IServer //当前zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称
	//zinx
	Version        string //当前zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包大小的最大值
}

/*
GlobalObject 定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

/*
Reload 从zinx.json去加载用于自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
Init 提供一个Init方法，初始化当前的GlobalObject
*/
func Init() {
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从conf/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
