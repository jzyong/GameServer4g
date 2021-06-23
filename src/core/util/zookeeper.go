package util

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"strings"
	"time"
)

const GateGameService string = "/mmo/%s/service/GateGameTcp/%d"     //gate game tcp服务
const GateClientService string = "/mmo/%s/service/GateClientTcp/%d" //gate client tcp服务
const HallRpcService string = "/mmo/%s/service/HallRpc/%d"          //hall rpc服务
const HallConfig string = "/mmo/%s/hall%d"                          //hall 配置
const GateConfig string = "/mmo/%s/gate%d"                          //gate 配置

//参考：https://www.cnblogs.com/zhichaoma/p/12640064.html
// https://blog.csdn.net/bingfeilongxin/article/details/88578887

//创建zookeeper连接
func ZKCreateConnect(hosts []string) *zk.Conn {
	connect, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		log.Error(err)
		return nil
	}
	return connect
}

// 增
// flags有4种取值：
// 0:永久，除非手动删除
// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
// zk.FlagSequence  = 2:会自动在节点后面添加序号
// 3:Ephemeral和Sequence，即，短暂且自动添加序号
func ZKAdd(conn *zk.Conn, path string, value string, flag int32) {
	var data = []byte(value)
	// 获取访问控制权限
	acls := zk.WorldACL(zk.PermAll)
	//先创建父节点
	paths := strings.Split(path, "/")
	var build strings.Builder
	for i, v := range paths {
		fmt.Sprintf(v)
		if i == 0 {
			build.WriteString("/")
			continue
		} else if i == len(paths)-1 {
			continue
		}
		build.WriteString(v)

		exist, _, err2 := conn.Exists(build.String())
		if err2 != nil {
			log.Warnf("创建节点%s %s失败，%v", path, build.String(), err2)
			return
		}
		if exist {
			build.WriteString("/")
			continue
		}
		//父节点永久存在
		conn.Create(build.String(), []byte(""), 0, acls)
		log.Infof("创建父节点：%s", build.String())
		build.WriteString("/")
	}
	s, err := conn.Create(path, data, flag, acls)
	if err != nil {
		log.Warnf("创建失败: %v\n", err)
		return
	}
	log.Infof("创建: %s 成功", s)
}

// 查
func ZKGet(conn *zk.Conn, path string) string {
	data, _, err := conn.Get(path)
	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", path, err)
		return ""
	}
	log.Infof("%s 的值为 %s\n", path, string(data))
	return string(data)
}

// 删改与增不同在于其函数中的version参数,其中version是用于 CAS支持
// 可以通过此种方式保证原子性
// 改
func ZKUpdate(conn *zk.Conn, path string, value string) {

	//先检查节点是否存在，不存在创建新的
	exist, _, err2 := conn.Exists(path)
	if err2 != nil {
		log.Warnf("更新节点%s 失败，%v", path, err2)
		return
	}
	if !exist {
		ZKAdd(conn, path, value, 0)
		return
	}

	newData := []byte(value)
	_, sate, _ := conn.Get(path)
	_, err := conn.Set(path, newData, sate.Version)
	if err != nil {
		fmt.Printf("数据修改失败: %v\n", err)
		return
	}
	log.Infof("%s 的值修改为 %s\n", path, value)
}

// 删
func ZKDelete(conn *zk.Conn, path string) {
	_, sate, _ := conn.Get(path)
	err := conn.Delete(path, sate.Version)
	if err != nil {
		fmt.Printf("数据删除失败: %v\n", err)
		return
	}
	log.Infof("路径%s 删除", path)
}
