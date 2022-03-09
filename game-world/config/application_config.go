package config

import (
	"encoding/json"
	"github.com/jzyong/golib/log"
	"io/ioutil"
	"os"
)

//配置
var ApplicationConfigInstance *ApplicationConfig

//配置文件路径
var FilePath string

//json配置对象
type ApplicationConfig struct {
	//服务器ID
	Id int32 `json:"id"`
	//rpc 地址
	RpcUrl string `json:"rpcUrl"`
	//日志级别
	LogLevel string "debug"
	//zookeeper 地址
	ZookeeperUrls []string `json:"zookeeperUrls"`
	//自定义配置
	Profile string `json:"profile"`
}

func init() {
	ApplicationConfigInstance = &ApplicationConfig{
		Id:       1,
		LogLevel: "debug",
	}
	//HallConfigInstance.Reload()
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
func (worldConfig *ApplicationConfig) Reload() {
	if confFileExists, _ := PathExists(FilePath); confFileExists != true {
		log.Warn("Config File ", FilePath, " is not exist!!")
		return
	}
	data, err := ioutil.ReadFile(FilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, worldConfig)
	if err != nil {
		log.Error("%v", err)
		panic(err)
	}
}
