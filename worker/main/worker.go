package main

import (
	"flag"
	"fmt"
	"github.com/c0ding/crontab/worker"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "worker,json")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	initArgs()
	initEnv()
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化任务管理器
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)

}
