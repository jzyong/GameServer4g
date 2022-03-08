module "github.com/jzyong/GameServer4g/game-api"

go 1.14

require (
	github.com/jzyong/GameServer4g/game-common v0.0.0 // indirect
	go.mongodb.org/mongo-driver v1.8.2
)

replace github.com/jzyong/GameServer4g/game-common => ../game-common