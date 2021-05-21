/**
*  针对 timer_wheel.go 时间轮api 做单元测试, 主要测试时间轮运转功能
*  依赖模块 runnable.go scheduled_task.go
 */
package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func TestTimerWheel(t *testing.T) {
	//创建秒级时间轮
	second_tw := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)

	//创建分钟级时间轮
	minute_tw := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)

	//创建小时级时间轮
	hour_tw := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	//将分层时间轮做关联
	hour_tw.AddTimeWheel(minute_tw)
	minute_tw.AddTimeWheel(second_tw)

	fmt.Println("init timeWheels done!")

	//===== > 以上为初始化分层时间轮 <====

	//给时间轮添加定时器
	timer1 := NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{1, 10}), 10*time.Second)
	_ = hour_tw.AddScheduledTask(1, timer1)
	fmt.Println("add timer 1 done!")

	//给时间轮添加定时器
	timer2 := NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{2, 20}), 20*time.Second)
	_ = hour_tw.AddScheduledTask(2, timer2)
	fmt.Println("add timer 2 done!")

	//给时间轮添加定时器
	timer3 := NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{3, 30}), 30*time.Second)
	_ = hour_tw.AddScheduledTask(3, timer3)
	fmt.Println("add timer 3 done!")

	//给时间轮添加定时器
	timer4 := NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{4, 40}), 40*time.Second)
	_ = hour_tw.AddScheduledTask(4, timer4)
	fmt.Println("add timer 4 done!")

	//给时间轮添加定时器
	timer5 := NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{5, 50}), 50*time.Second)
	_ = hour_tw.AddScheduledTask(5, timer5)
	fmt.Println("add timer 5 done!")

	//时间轮运行
	second_tw.Run()
	minute_tw.Run()
	hour_tw.Run()

	fmt.Println("timeWheels are run!")

	go func() {
		n := 0.0
		for {
			fmt.Println("tick...", n)

			//取出近1ms的超时定时器有哪些
			tasks := hour_tw.GetScheduledTaskWithIn(1000 * time.Millisecond)
			for _, task := range tasks {
				//调用定时器方法
				task.runnable.Run()
			}

			time.Sleep(500 * time.Millisecond)
			n += 0.5
		}
	}()

	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	time.Sleep(10 * time.Minute)
}
