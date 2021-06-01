package network

type HandlerMethod func(message TcpMessage) bool

/*
	路由接口， 这里面路由是 使用框架者给该链接自定的 处理业务方法
	路由里的TcpMessage 则包含用该链接的链接信息和该链接的请求数据信息
*/
type TcpHandler struct {
	run HandlerMethod //处理的方法
}

func NewTcpHandler(method HandlerMethod) *TcpHandler {
	return &TcpHandler{
		run: method,
	}
}

//获取运行方法
func (this *TcpHandler) GetRunMethod() HandlerMethod {
	return this.run
}
