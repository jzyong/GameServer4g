package concurrent

import (
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"reflect"
	"time"
)

const (
	HOUR_NAME     = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3 //ms为精度
	HOUR_SCALES   = 12

	MINUTE_NAME     = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES   = 60

	SECOND_NAME     = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES   = 60

	TIMERS_MAX_CAP = 2048 //每个时间轮刻度挂载定时器的最大个数
)

/*
	时间任务
*/
type ScheduledTask struct {
	//延迟调用函数
	runnable *Runnable
	//调用时间(unix 时间， 单位ms)
	unixTime int64
}

//返回1970-1-1至今经历的毫秒数
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

/*
   创建一个定时器,在指定的时间触发 定时器方法
	runnable: Runnable类型的延迟调用函数类型
	unixNano: unix计算机从1970-1-1至今经历的纳秒数
*/
func NewScheduledTask(runnable *Runnable, unixNano int64) *ScheduledTask {
	return &ScheduledTask{
		runnable: runnable,
		unixTime: unixNano / 1e6, //将纳秒转换成对应的毫秒 ms ，定时器以ms为最小精度
	}
}

/*
	创建一个定时器，在当前时间延迟duration之后触发 定时器方法
*/
func NewScheduledTaskAfter(runnable *Runnable, duration time.Duration) *ScheduledTask {
	return NewScheduledTask(runnable, time.Now().UnixNano()+int64(duration))
}

//启动定时器，用一个go承载
func (s *ScheduledTask) Run() {
	go func() {
		now := UnixMilli()
		//设置的定时器是否在当前时间之后
		if s.unixTime > now {
			//睡眠，直至时间超时,已微秒为单位进行睡眠
			time.Sleep(time.Duration(s.unixTime-now) * time.Millisecond)
		}
		//调用事先注册好的超时延迟方法
		s.runnable.Run()
	}()
}

/*
   定义一个延迟调用函数
	延迟调用函数就是 时间定时器超时的时候，触发的事先注册好的
	回调函数
*/
type Runnable struct {
	run  func(...interface{}) //run : 延迟函数调用原型
	args []interface{}        //args: 延迟调用函数传递的形参
}

/*
	创建一个延迟调用函数
*/
func NewRunnable(run func(v ...interface{}), args []interface{}) *Runnable {
	return &Runnable{
		run:  run,
		args: args,
	}
}

//打印当前延迟函数的信息，用于日志记录
func (r *Runnable) String() string {
	return fmt.Sprintf("{Runnable:%s, args:%v}", reflect.TypeOf(r.run).Name(), r.args)
}

/*
	执行延迟函数---如果执行失败，抛出异常
*/
func (r *Runnable) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(r.String(), "Run err: ", err)
		}
	}()

	//调用定时器超时函数
	r.run(r.args...)
}
