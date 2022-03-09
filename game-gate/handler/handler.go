package handler

import (
	"github.com/jzyong/GameServer4g/game-gate/manager"
	"github.com/jzyong/GameServer4g/game-message/message"
)

//注册client handler
func RegisterClientHandler() {
	manager.GetClientManager().RegisterHandler(int32(message.MID_PlayerHeartReq), HandlePlayerHeartReq)
}

//注册game handler
func RegisterGameHandler() {
	manager.GetGameManager().RegisterHandler(int32(message.MID_ServerRegisterUpdateReq), HandleServerRegister)
}
