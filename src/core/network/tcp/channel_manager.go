package network

import (
	"errors"
	"fmt"
	"sync"
)

/*
	连接管理抽象层
*/
type ChannelManager interface {
	Add(conn Channel)                   //添加链接
	Remove(conn Channel)                //删除连接
	Get(connID uint32) (Channel, error) //利用ConnID获取链接
	Len() int                           //获取当前连接
	ClearConn()                         //删除并停止所有链接
}

/*
	连接管理模块
*/
type channelManagerImpl struct {
	channels map[uint32]Channel //管理的连接信息
	connLock sync.RWMutex       //读写连接的读写锁
}

/*
	创建一个链接管理
*/
func NewChannelManager() ChannelManager {
	return &channelManagerImpl{
		channels: make(map[uint32]Channel),
	}
}

//添加链接
func (connMgr *channelManagerImpl) Add(conn Channel) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//将conn连接添加到ConnMananger中
	connMgr.channels[conn.GetConnID()] = conn

	//log.Debug("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

//删除连接
func (connMgr *channelManagerImpl) Remove(conn Channel) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.channels, conn.GetConnID())

	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
}

//利用ConnID获取链接
func (connMgr *channelManagerImpl) Get(connID uint32) (Channel, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.channels[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

//获取当前连接数
func (connMgr *channelManagerImpl) Len() int {
	return len(connMgr.channels)
}

//清除并停止所有连接
func (connMgr *channelManagerImpl) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.channels {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.channels, connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}
