package rpc

/**
rpc 客户端接口
*/
type IGrpcClient interface {
	//启动
	Start(rpcUrl string)

	//停止
	Stop()

	////连接
	//Connect() grpc.ClientConnInterface
}
