/**
*  时间轮调度器
*   依赖模块 scheduled_task.go time_wheel.go
 */
package concurrent

import (
	"github.com/jzyong/go-mmo-server/src/core/log"
	"math"
	"sync"
	"time"
)

const (
	//默认缓冲触发函数队列大小
	MAX_CHAN_BUFF = 2048
	//默认最大误差时间
	MAX_TIME_DELAY = 100
)

//调度池
type ScheduledPool struct {
	//当前调度器的最高级时间轮
	tw *TimeWheel
	//定时器编号累加器
	idGen uint32
	//已经触发定时器的channel
	triggerChan chan *Runnable
	//互斥锁
	sync.RWMutex
}

/*
	返回一个定时器调度器

	主要创建分层定时器，并做关联，并依次启动
*/
func NewScheduledPool() *ScheduledPool {

	//创建秒级时间轮
	second_tw := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	//创建分钟级时间轮
	minute_tw := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	//创建小时级时间轮
	hour_tw := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	//将分层时间轮做关联
	hour_tw.AddTimeWheel(minute_tw)
	minute_tw.AddTimeWheel(second_tw)

	//时间轮运行
	second_tw.Run()
	minute_tw.Run()
	hour_tw.Run()

	return &ScheduledPool{
		tw:          hour_tw,
		triggerChan: make(chan *Runnable, MAX_CHAN_BUFF),
	}
}

//创建一个定点Timer 并将Timer添加到分层时间轮中， 返回Timer的tid
func (this *ScheduledPool) CreateTimerAt(runnable *Runnable, unixNano int64) (uint32, error) {
	this.Lock()
	defer this.Unlock()
	this.idGen++
	return this.idGen, this.tw.AddScheduledTask(this.idGen, NewScheduledTask(runnable, unixNano))
}

//创建一个延迟Timer 并将Timer添加到分层时间轮中， 返回Timer的tid
func (this *ScheduledPool) CreateTimerAfter(df *Runnable, duration time.Duration) (uint32, error) {
	this.Lock()
	defer this.Unlock()
	this.idGen++
	return this.idGen, this.tw.AddScheduledTask(this.idGen, NewScheduledTaskAfter(df, duration))
}

//删除timer
func (this *ScheduledPool) CancelTimer(id uint32) {
	this.Lock()
	this.Unlock()

	this.tw.RemoveScheduledTask(id)
}

//获取计时结束的延迟执行函数通道
func (this *ScheduledPool) GetTriggerChan() chan *Runnable {
	return this.triggerChan
}

//非阻塞的方式启动timerSchedule，只取出，不执行任务
func (this *ScheduledPool) Start() {
	go func() {
		for {
			//当前时间
			now := UnixMilli()
			//获取最近MAX_TIME_DELAY 毫秒的超时定时器集合
			tasks := this.tw.GetScheduledTaskWithIn(MAX_TIME_DELAY * time.Millisecond)
			for _, task := range tasks {
				if math.Abs(float64(now-task.unixTime)) > MAX_TIME_DELAY {
					//已经超时的定时器，报警
					log.Error("want call at ", task.unixTime, "; real call at", now, "; delay ", now-task.unixTime)
				}
				//将超时触发函数写入管道
				this.triggerChan <- task.runnable
			}
			time.Sleep(MAX_TIME_DELAY / 2 * time.Millisecond)
		}
	}()
}

//时间轮定时器 自动调度
func NewAutoExecuteScheduledPool() *ScheduledPool {
	//创建一个调度器
	autoExecuteScheduler := NewScheduledPool()
	//启动调度器
	autoExecuteScheduler.Start()

	//永久从调度器中获取超时 触发的函数 并执行
	go func() {
		delayFuncChan := autoExecuteScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			go df.Run()
		}
	}()

	return autoExecuteScheduler
}
