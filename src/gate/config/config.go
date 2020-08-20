package config

import (
	"encoding/json"
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"io/ioutil"
	"os"
)

//网关配置
var GateConfigInstance *GateConfig

//网关json配置对象
type GateConfig struct {
	//服务器ID
	Id int32
	//允许用户连接数
	UserConnectCount int32
	//配置文件路径
	ConfigFilePath string "gate"
	//日志级别
	LogLevel string "debug"
	//日志名称
	LogFileName string
	//gate rpc地址
	ClusterRpcURL string
}

func init() {
	GateConfigInstance = &GateConfig{
		Id:               2,
		LogLevel:         "debug",
		ConfigFilePath:   "config/GateConfig.json",
		UserConnectCount: 10000,
		ClusterRpcURL:    "192.168.110.16:2002",
	}
	GateConfigInstance.Reload()
}

//判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//读取用户的配置文件
func (g *GateConfig) Reload() {
	if confFileExists, _ := PathExists(g.ConfigFilePath); confFileExists != true {
		fmt.Println("Config File ", g.ConfigFilePath, " is not exist!!")
		return
	}
	data, err := ioutil.ReadFile(g.ConfigFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
