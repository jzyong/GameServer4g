package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/gate/manager"
	"github.com/jzyong/go-mmo-server/src/message"
)

//获取服务器列表 (遗弃)
func HandleServerList(msg network.TcpMessage) bool {
	request := &message.ServerListRequest{}
	proto.Unmarshal(msg.GetData(), request)
	log.Infof("请求%d", request.GetType())

	var serverInfo []*message.ServerInfo
	serverInfo = []*message.ServerInfo{
		{
			Id: 1,
			Ip: "",
		},
	}
	response := &message.ServerListResponse{
		Server: serverInfo,
	}

	network.SendClientProtoMsg(msg.GetChannel(), int32(message.MID_ServerListRes), response)
	return true
}

//后端服务器注册
func HandleServerRegister(msg network.TcpMessage) bool {
	request := &message.ServerRegisterUpdateRequest{}
	proto.Unmarshal(msg.GetData(), request)
	serverInfo := request.GetServerInfo()
	//log.Infof("server %d %d: %s register to gate state %d", serverInfo.GetId(), serverInfo.GetType(), serverInfo.GetIp(), serverInfo.GetState())

	manager.GetGameManager().UpdateHallServerInfo(serverInfo, msg.GetChannel())
	//TODO 加到连接管理中

	return true
}
