package network

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"strconv"
)

/*
	消息分发管理抽象层
*/
type MessageDistribute interface {
	//处理消息
	RunHandler(message TcpMessage)
	//为消息添加具体的处理逻辑
	RegisterHandler(msgId int32, handler *TcpHandler)
	//启动worker工作池
	StartWorkerPool()
	//将消息交给TaskQueue,由worker进行处理
	SendMessageToTaskQueue(message TcpMessage)
}

//处理未注册消息，如转发到大厅
type HandUnregisterMessageMethod func(message TcpMessage)

//Handler 处理器
type messageDistributeImpl struct {
	handlers              map[int32]*TcpHandler       //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize        uint32                      //业务工作Worker池的数量
	TaskQueue             []chan TcpMessage           //Worker负责取任务的消息队列
	HandUnregisterMessage HandUnregisterMessageMethod //处理未注册消息，如转发到大厅
}

func NewMessageDistribute(workPoolSize uint32, unregisterMethod HandUnregisterMessageMethod) MessageDistribute {
	return &messageDistributeImpl{
		handlers:       make(map[int32]*TcpHandler),
		WorkerPoolSize: workPoolSize,
		//一个worker对应一个queue
		TaskQueue:             make([]chan TcpMessage, workPoolSize),
		HandUnregisterMessage: unregisterMethod,
	}
}

//将消息交给TaskQueue,由worker进行处理
func (mh *messageDistributeImpl) SendMessageToTaskQueue(request TcpMessage) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetChannel().GetConnID() % mh.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

//马上以非阻塞方式处理消息
func (mh *messageDistributeImpl) RunHandler(msg TcpMessage) {
	handler, ok := mh.handlers[msg.GetMsgId()]
	if !ok {
		if mh.HandUnregisterMessage != nil {
			mh.HandUnregisterMessage(msg)
			return
		}
		log.Warn("Handler msgId = ", msg.GetMsgId(), " is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.run(msg)
}

//为消息添加具体的处理逻辑
func (mh *messageDistributeImpl) RegisterHandler(msgId int32, handler *TcpHandler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.handlers[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与handler的绑定关系
	mh.handlers[msgId] = handler
	log.Infof("Add handler %d ", msgId)
}

//启动一个Worker工作流程
func (mh *messageDistributeImpl) StartOneWorker(workerID int, taskQueue chan TcpMessage) {
	log.Info("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.RunHandler(request)
		}
	}
}

//启动worker工作池
func (mh *messageDistributeImpl) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan TcpMessage, 1024)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
