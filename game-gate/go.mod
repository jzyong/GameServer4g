module github.com/jzyong/GameServer4g/game-gate

go 1.19

require (
	github.com/go-zookeeper/zk v1.0.2
	github.com/golang/protobuf v1.5.2
	github.com/jzyong/GameServer4g/game-common v0.0.0-00010101000000-000000000000
	github.com/jzyong/GameServer4g/game-message v0.0.0-00010101000000-000000000000
	github.com/jzyong/golib v0.0.16
)

require (
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.44.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)

replace (
	github.com/jzyong/GameServer4g/game-common => ../game-common
	github.com/jzyong/GameServer4g/game-message => ../game-message
)
