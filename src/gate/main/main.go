package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/jzy/go-mmo-server/src/message"
	"log"
)

func main()  {
	testProto()
}

//测试protobuf
func testProto()  {
	serverInfo :=&message.ServerInfo{
		Id:                   1,
		BelongID:             1,
		Ip:                   "192.168.0.1",
		Type:                 1,
		Port:                 0,
		State:                0,
		Version:              "",
		Content:              "",
		Online:               1000,
		MaxUserCount:         1000,
		HttpPort:             8080,
		OpenTime:             "",
		MaintainTime:         "",
		Name:                 "",
		Wwwip:                "",
	}
	log.Println("server info :",serverInfo.String())
	data,err:=proto.Marshal(serverInfo)
	if err!=nil{
		log.Fatal("marshaling error:",err)
	}
	newServerInfo:=&message.ServerInfo{}
	err=proto.Unmarshal(data,newServerInfo)
	if err!=nil{
		log.Fatal("unmarshaling error:",err)
	}
	if serverInfo.GetId()!=newServerInfo.GetId(){
		log.Fatal("data mismatch %v != %v",serverInfo.GetId(),newServerInfo.GetId())
	}
}