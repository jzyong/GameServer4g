package handler

import (
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/core/log"
	network "github.com/jzyong/go-mmo-server/src/core/network/tcp"
	"github.com/jzyong/go-mmo-server/src/gate/manager"
	"github.com/jzyong/go-mmo-server/src/message"
)

func init() {
	manager.GetClientManager().GetServer().RegisterHandler(int32(message.MID_ServerListReq), &HelloHandler{})
}

type HelloHandler struct {
	network.BaseTcpHandler
}

func (br *HelloHandler) Run(msg network.TcpMessage) {
	request := &message.ServerListRequest{}
	proto.Unmarshal(msg.GetData(), request)
	log.Infof("请求%d", request.GetType())

	var serverInfo []*message.ServerInfo
	serverInfo = []*message.ServerInfo{
		{
			Id:   1,
			Ip:   "",
			Name: "111",
		},
	}
	response := &message.ServerListResponse{
		Server: serverInfo,
	}

	network.SendClientProtoMsg(msg.GetChannel(), int32(message.MID_ServerListRes), response)

}
