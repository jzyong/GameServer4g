package concurrent

import (
	"errors"
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"sync"
	"time"
)

/*
  tips:
	一个网络服务程序时需要管理大量客户端连接的，
	其中每个客户端连接都需要管理它的 timeout 时间。
	通常连接的超时管理一般设置为30~60秒不等，并不需要太精确的时间控制。
	另外由于服务端管理着多达数万到数十万不等的连接数，
	因此我们没法为每个连接使用一个Timer，那样太消耗资源不现实。

	用时间轮的方式来管理和维护大量的timer调度，会解决上面的问题。
*/
type TimeWheel struct {
	//TimeWheel的名称
	name string
	//刻度的时间间隔，单位ms
	interval int64
	//每个时间轮上的刻度数
	scales int
	//当前时间指针的指向
	curIndex int
	//每个刻度所存放的timer定时器的最大容量
	maxCap int
	//当前时间轮上的所有timer
	scheduledTaskQueue map[int]map[uint32]*ScheduledTask //map[int] VALUE  其中int表示当前时间轮的刻度,map[int] map[uint32] *ScheduledTask, uint32表示Timer的id号
	//下一层时间轮
	nextTimeWheel *TimeWheel
	//互斥锁（继承RWMutex的 RWLock,UnLock 等方法）
	sync.RWMutex
}

/*
	创建一个时间轮
	name：时间轮的名称
	interval：每个刻度之间的duration时间间隔
	scales:当前时间轮的轮盘一共多少个刻度(如我们正常的时钟就是12个刻度)
	maxCap: 每个刻度所最大保存的Timer定时器个数
*/
func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	tw := &TimeWheel{
		name:               name,
		interval:           interval,
		scales:             scales,
		maxCap:             maxCap,
		scheduledTaskQueue: make(map[int]map[uint32]*ScheduledTask, scales),
	}
	//初始化map
	for i := 0; i < scales; i++ {
		tw.scheduledTaskQueue[i] = make(map[uint32]*ScheduledTask, maxCap)
	}
	log.Info("Init timerWheel name = ", tw.name, " is Done!")
	return tw
}

/*
	将一个ScheduledTask定时器加入到分层时间轮中
	id: 每个定时器ScheduledTask的唯一标识
	t: 当前被加入时间轮的定时器
	forceNext: 是否强制的将定时器添加到下一层时间轮

	我们采用的算法是：
	如果当前ScheduledTask的超时时间间隔 大于一个刻度，那么进行hash计算 找到对应的刻度上添加
	如果当前的ScheduledTask的超时时间间隔 小于一个刻度 :
					如果没有下一轮时间轮
*/
func (tw *TimeWheel) addScheduledTask(id uint32, scheduledTask *ScheduledTask, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errStr := fmt.Sprintf("addScheduledTask function err : %s", err)
			log.Error(errStr)
			return errors.New(errStr)
		}
		return nil
	}()

	//得到当前的超时时间间隔(ms)毫秒为单位
	delayInterval := scheduledTask.unixTime - UnixMilli()

	//如果当前的超时时间 大于一个刻度的时间间隔
	if delayInterval >= tw.interval {
		//得到需要跨越几个刻度
		dn := delayInterval / tw.interval
		//在对应的刻度上的定时器Timer集合map加入当前定时器(由于是环形，所以要求余)
		tw.scheduledTaskQueue[(tw.curIndex+int(dn))%tw.scales][id] = scheduledTask
		return nil
	}

	//如果当前的超时时间,小于一个刻度的时间间隔，并且当前时间轮没有下一层，经度最小的时间轮
	if delayInterval < tw.interval && tw.nextTimeWheel == nil {
		if forceNext == true {
			//如果设置为强制移至下一个刻度，那么将定时器移至下一个刻度
			//这种情况，主要是时间轮自动轮转的情况
			//因为这是底层时间轮，该定时器在转动的时候，如果没有被调度者取走的话，该定时器将不会再被发现
			//因为时间轮刻度已经过去，如果不强制把该定时器scheduledTask移至下时刻，就永远不会被取走并触发调用
			//所以这里强制将scheduledTask移至下个刻度的集合中，等待调用者在下次轮转之前取走该定时器
			tw.scheduledTaskQueue[(tw.curIndex+1)%tw.scales][id] = scheduledTask
		} else {
			//如果手动添加定时器，那么直接将scheduledTask添加到对应底层时间轮的当前刻度集合中
			tw.scheduledTaskQueue[tw.curIndex][id] = scheduledTask
		}
		return nil
	}

	//如果当前的超时时间，小于一个刻度的时间间隔，并且有下一层时间轮
	if delayInterval < tw.interval {
		return tw.nextTimeWheel.AddScheduledTask(id, scheduledTask)
	}

	return nil
}

//添加一个ScheduledTask到一个时间轮中(非时间轮自转情况)
func (tw *TimeWheel) AddScheduledTask(id uint32, scheduledTask *ScheduledTask) error {
	tw.Lock()
	defer tw.Unlock()

	return tw.addScheduledTask(id, scheduledTask, false)
}

/*
	删除一个定时器，根据定时器的id
*/
func (tw *TimeWheel) RemoveScheduledTask(id uint32) {
	tw.Lock()
	defer tw.Unlock()

	for i := 0; i < tw.scales; i++ {
		if _, ok := tw.scheduledTaskQueue[i][id]; ok {
			delete(tw.scheduledTaskQueue[i], id)
		}
	}
}

/*
	给一个时间轮添加下层时间轮 比如给小时时间轮添加分钟时间轮，给分钟时间轮添加秒时间轮
*/
func (tw *TimeWheel) AddTimeWheel(next *TimeWheel) {
	tw.nextTimeWheel = next
	log.Info("Add timerWhell[", tw.name, "]'s next [", next.name, "] is success!")
}

/*
	启动时间轮
*/
func (tw *TimeWheel) run() {
	for {
		//时间轮每间隔interval一刻度时间，触发转动一次
		time.Sleep(time.Duration(tw.interval) * time.Millisecond)

		tw.Lock()
		//取出挂载在当前刻度的全部定时器
		curScheduledTasks := tw.scheduledTaskQueue[tw.curIndex]
		//当前定时器要重新添加 所给当前刻度再重新开辟一个map ScheduledTask容器
		tw.scheduledTaskQueue[tw.curIndex] = make(map[uint32]*ScheduledTask, tw.maxCap)
		for id, scheduledTask := range curScheduledTasks {
			//这里属于时间轮自动转动，forceNext设置为true
			tw.addScheduledTask(id, scheduledTask, true)
		}

		//取出下一个刻度 挂载的全部定时任务 进行重新添加 (为了安全起见,待考慮)
		nextScheduledTasks := tw.scheduledTaskQueue[(tw.curIndex+1)%tw.scales]
		tw.scheduledTaskQueue[(tw.curIndex+1)%tw.scales] = make(map[uint32]*ScheduledTask, tw.maxCap)
		for id, scheduledTask := range nextScheduledTasks {
			tw.addScheduledTask(id, scheduledTask, true)
		}

		//当前刻度指针 走一格
		tw.curIndex = (tw.curIndex + 1) % tw.scales

		tw.Unlock()
	}
}

//非阻塞的方式让时间轮转起来
func (tw *TimeWheel) Run() {
	go tw.run()
	log.Info("timerWheel name = ", tw.name, " is running...")
}

//获取定时器在一段时间间隔内的Timer
func (tw *TimeWheel) GetScheduledTaskWithIn(duration time.Duration) map[uint32]*ScheduledTask {
	//最终触发定时器的一定是挂载最底层时间轮上的定时器
	//1 找到最底层时间轮
	leaftw := tw
	for leaftw.nextTimeWheel != nil {
		leaftw = leaftw.nextTimeWheel
	}

	leaftw.Lock()
	defer leaftw.Unlock()
	//返回的Timer集合
	taskList := make(map[uint32]*ScheduledTask)

	now := UnixMilli()

	//取出当前时间轮刻度内全部Timer
	for id, scheduledTask := range leaftw.scheduledTaskQueue[leaftw.curIndex] {
		if scheduledTask.unixTime-now < int64(duration/1e6) {
			//当前定时器已经超时
			taskList[id] = scheduledTask
			//定时器已经超时被取走，从当前时间轮上 摘除该定时器
			delete(leaftw.scheduledTaskQueue[leaftw.curIndex], id)
		}
	}
	return taskList
}
