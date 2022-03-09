module github.com/jzyong/GameServer4g/game-res

go 1.14

require (
	github.com/jzyong/GameServer4g/game-message v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.44.0
)
replace (
	github.com/jzyong/GameServer4g/game-message => ../game-message
)