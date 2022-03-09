package handler

import (
	"github.com/jzyong/GameServer4g/game-hall/manager"
	"github.com/jzyong/GameServer4g/game-message/message"
)

//注册消息处理器
func RegisterHandlers() {
	manager.GetClientManager().RegisterHandler(message.MID_UserLoginReq, HandUserLogin)
}
