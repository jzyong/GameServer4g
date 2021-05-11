package network

/*
	路由接口， 这里面路由是 使用框架者给该链接自定的 处理业务方法
	路由里的TcpMessage 则包含用该链接的链接信息和该链接的请求数据信息
*/
type TcpHandler interface {
	//在处理conn业务之前的钩子方法
	PreRun(message TcpMessage)
	//处理conn业务的方法
	Run(message TcpMessage)
	//处理conn业务之后的钩子方法
	PostRun(message TcpMessage)
}

//实现TcpHandler时，先嵌入这个基类，然后根据需要对这个基类的方法进行重写
type BaseTcpHandler struct{}

//这里之所以BaseRouter的方法都为空，
// 是因为有的Router不希望有PreHandle或PostHandle
// 所以Router全部继承BaseRouter的好处是，不需要实现PreHandle和PostHandle也可以实例化
func (br *BaseTcpHandler) PreRun(message TcpMessage)  {}
func (br *BaseTcpHandler) Run(message TcpMessage)     {}
func (br *BaseTcpHandler) PostRun(message TcpMessage) {}
