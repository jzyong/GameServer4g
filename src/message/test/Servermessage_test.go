package test

import (
	"github.com/golang/protobuf/proto"
	"github.com/jzyong/go-mmo-server/src/message"
	"log"
	"testing"
)

//serverInfo 构建
func TestServerInfo(t *testing.T) {
	t.Run("test serverInfo", func(t *testing.T) {
		buildServerInfo()
	})
}

//测试protobuf
func buildServerInfo() {
	serverInfo := &message.ServerInfo{
		Id:           1,
		BelongID:     1,
		Ip:           "192.168.0.1",
		Type:         1,
		Port:         0,
		State:        0,
		Version:      "",
		Content:      "",
		Online:       1000,
		MaxUserCount: 1000,
		HttpPort:     8080,
		OpenTime:     "",
		MaintainTime: "",
		Name:         "",
		Wwwip:        "",
	}
	log.Println("server info :", serverInfo.String())
	data, err := proto.Marshal(serverInfo)
	if err != nil {
		log.Fatal("marshaling error:", err)
	}
	newServerInfo := &message.ServerInfo{}
	err = proto.Unmarshal(data, newServerInfo)
	if err != nil {
		log.Fatal("unmarshaling error:", err)
	}
	if serverInfo.GetId() != newServerInfo.GetId() {
		log.Fatal("data mismatch ", serverInfo.GetId(), newServerInfo.GetId())
	}
}
