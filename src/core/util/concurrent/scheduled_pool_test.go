/**
*  时间轮定时器调度器单元测试
 */
package concurrent

import (
	"fmt"
	"github.com/jzyong/go-mmo-server/src/core/log"
	"testing"
	"time"
)

//触发函数
func foo(args ...interface{}) {
	fmt.Printf("I am No. %d function, delay %d ms\n", args[0].(int), args[1].(int))
}

//手动创建调度器运转时间轮
func TestNewScheduledPool(t *testing.T) {
	timerScheduler := NewScheduledPool()
	timerScheduler.Start()

	//在scheduler中添加task
	for i := 1; i < 2000; i++ {
		f := NewRunnable(foo, []interface{}{i, i * 3})
		tid, err := timerScheduler.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Error("create timer error", tid, err)
			break
		}
	}

	//执行调度器触发函数
	go func() {
		runnables := timerScheduler.GetTriggerChan()
		for runnable := range runnables {
			runnable.Run()
		}
	}()

	//阻塞等待
	select {}
}

//采用自动调度器运转时间轮
func TestNewAutoExecuteScheduledPool(t *testing.T) {
	autoTS := NewAutoExecuteScheduledPool()

	//给调度器添加task
	for i := 0; i < 2000; i++ {
		f := NewRunnable(foo, []interface{}{i, i * 3})
		tid, err := autoTS.CreateTimerAfter(f, time.Duration(3*i)*time.Millisecond)
		if err != nil {
			log.Error("create timer error", tid, err)
			break
		}
	}

	//阻塞等待
	select {}
}
