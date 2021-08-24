package handler

import (
	"github.com/jzyong/go-mmo-server/src/hall/manager"
	"github.com/jzyong/go-mmo-server/src/message"
)

//注册消息处理器
func RegisterHandlers() {
	manager.GetClientManager().RegisterHandler(message.MID_UserLoginReq, HandUserLogin)
}
