module github.com/jzyong/GameServer4g/game-gate

go 1.14

require (
	github.com/go-zookeeper/zk v1.0.2
	github.com/golang/protobuf v1.5.2
	github.com/jzyong/GameServer4g/game-message v0.0.0-00010101000000-000000000000
	github.com/jzyong/GameServer4g/game-common v0.0.0-00010101000000-000000000000
	github.com/jzyong/golib v0.0.8
)

replace (
	github.com/jzyong/GameServer4g/game-common => ../game-common
	github.com/jzyong/GameServer4g/game-message => ../game-message
)
