package game_common

const (
	GateGameService           string = "/game/%s/service/GateGameTcp/%d"   //gate game tcp服务
	GateGameServiceListenPath string = "/game/%s/service/GateGameTcp"      //网关服务监听路径
	GateClientService         string = "/game/%s/service/GateClientTcp/%d" //gate client tcp服务
	HallRpcService            string = "/game/%s/service/HallRpc/%d"       //hall rpc服务
	WorldRpcService           string = "/game/%s/service/WorldRpc/%d"      //world rpc服务
	WorldRpcServiceListenPath string = "/game/%s/service/WorldRpc"         //world rpc服务
	HallConfig                string = "/game/%s/hall%d"                   //hall 配置
	GateConfig                string = "/game/%s/gate%d"                   //gate 配置
	WorldConfig               string = "/game/%s/world%d"                  //world 配置
)
