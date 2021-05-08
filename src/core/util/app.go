package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// 模块化
type Module interface {
	Init() error
	Run()
	Stop()
}

//@
type DefaultModule struct {
}

func (this DefaultModule) Init() error {
	return nil
}

func (this DefaultModule) Run() {

}

func (this DefaultModule) Stop() {

}

//  DefaultModuleManager default module manager
type DefaultModuleManager struct {
	Module
	Modules []Module
}

//
func NewDefaultModuleManager() *DefaultModuleManager {
	return &DefaultModuleManager{
		Modules: make([]Module, 0, 5),
	}
}

// 初始化所有模块
func (this *DefaultModuleManager) Init() error {
	for i := 0; i < len(this.Modules); i++ {
		err := this.Modules[i].Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// 运行模块
func (this *DefaultModuleManager) Run() {
	for i := 0; i < len(this.Modules); i++ {
		this.Modules[i].Run()
	}
}

func (this *DefaultModuleManager) Stop() {
	var wg sync.WaitGroup
	for i := 0; i < len(this.Modules); i++ {
		wg.Add(1)
		go func(module Module) {
			module.Stop()
			wg.Done()
		}(this.Modules[i])
	}
	wg.Wait()
}

// 添加模块
func (this *DefaultModuleManager) AppendModule(module Module) Module {
	this.Modules = append(this.Modules, module)
	return module
}

//  WaitTerminateSignal wait signal to end the program
func WaitForTerminate() {
	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		close(exitChan)
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
}
