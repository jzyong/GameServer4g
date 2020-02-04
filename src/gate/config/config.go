package config

import (
	"core/log"
	"encoding/json"
	"io/ioutil"
	innerLog "log"
)

//网关全局配置

//定义常量，变量
var ()

//读取配置文件
func init() {
	data, err := ioutil.ReadFile("config/gateConfig.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &GateConfig)
	if err != nil {
		log.Fatal("%v", err)
	}

	//设置日志
	if GateConfig.LogFileName != "" {
		defaultLog, err := log.New(GateConfig.LogLevel, "log", GateConfig.LogFileName, innerLog.LstdFlags)
		if err != nil {
			log.Fatal("%v", err)
		}
		log.DefaultLogger = defaultLog
	}

	log.Info("server config: %v", GateConfig)
}

//网关json配置对象
var GateConfig struct {
	//服务器ID
	Id int32
	//允许用户连接数
	UserConnectCount int32

	//日志文件名字,不设置控制台输出
	LogFileName string "gate"
	//日志级别
	LogLevel string "debug"
}
