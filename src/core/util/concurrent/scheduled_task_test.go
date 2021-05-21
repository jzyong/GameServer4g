/**
* 针对timer.go做单元测试，主要测试定时器相关接口 依赖模块delayFunc.go
 */
package concurrent

import (
	"fmt"
	"testing"
	"time"
)

func SayHello(message ...interface{}) {
	fmt.Println(message[0].(string), " ", message[1].(string))
}

func TestRunnable(t *testing.T) {
	runnable := NewRunnable(SayHello, []interface{}{"hello", "zinx!"})
	fmt.Println("runnable.String() = ", runnable.String())
	runnable.Run()
}

//定义一个超时函数
func myFunc(v ...interface{}) {
	fmt.Printf("No.%d function calld. delay %d second(s)\n", v[0].(int), v[1].(int))
}

func TestScheduledTask(t *testing.T) {

	for i := 0; i < 5; i++ {
		go func(i int) {
			NewScheduledTaskAfter(NewRunnable(myFunc, []interface{}{i, 2 * i}), time.Duration(2*i)*time.Second).Run()
		}(i)
	}

	//主进程等待其他go，由于Run()方法是用一个新的go承载延迟方法，这里不能用waitGroup
	time.Sleep(1 * time.Minute)
}
