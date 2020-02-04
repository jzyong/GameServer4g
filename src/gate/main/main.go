package main

import (
	"gate/config"
	"log"
)

func main() {
	log.Printf("gate:%d starting", config.GateConfig.Id)

}
