package config

import (
	"encoding/json"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"io/ioutil"
	"os"
)

//网关配置
var GateConfigInstance *GateConfig

//配置文件路径
var FilePath string

//网关json配置对象
type GateConfig struct {
	//服务器ID
	Id int32 `json:"id"`
	//客户端链接地址
	ClientUrl string `json:"clientUrl"`
	//后端游戏服务器地址
	GameUrl string `json:"gameUrl"`

	//允许用户连接数
	UserConnectCount int32 `json:"userConnectCount"`
	//日志级别
	LogLevel string "debug"
	//日志名称
	LogFileName string
	//gate rpc地址
	ClusterRpcURL string
	//zookeeper 地址
	ZookeeperUrls []string `json:"zookeeperUrls"`
	//自定义配置
	Profile string `json:"profile"`
}

func init() {
	GateConfigInstance = &GateConfig{
		Id:               2,
		ClientUrl:        "127.0.0.1:6060",
		GameUrl:          "127.0.0.1:6061",
		LogLevel:         "debug",
		UserConnectCount: 10000,
		ClusterRpcURL:    "192.168.110.16:2002",
	}
	//GateConfigInstance.Reload()
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
	if confFileExists, _ := PathExists(FilePath); confFileExists != true {
		log.Warn("Config File ", FilePath, " is not exist!!")
		return
	}
	data, err := ioutil.ReadFile(FilePath)
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
