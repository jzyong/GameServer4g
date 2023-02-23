module github.com/jzyong/GameServer4g/game-world

go 1.19

require (
	github.com/go-zookeeper/zk v1.0.2
	github.com/jzyong/GameServer4g/game-common v0.0.0-00010101000000-000000000000
	github.com/jzyong/GameServer4g/game-message v0.0.0-00010101000000-000000000000
	github.com/jzyong/golib v0.0.16
	google.golang.org/grpc v1.44.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.3.8 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace (
	github.com/jzyong/GameServer4g/game-common => ../game-common
	github.com/jzyong/GameServer4g/game-message => ../game-message
)
