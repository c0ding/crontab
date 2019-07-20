package main

import (
	"flag"
	"fmt"
	"github.com/c0ding/crontab/master"
	"runtime"
	"time"
)

var (
	confFile string // 通过解析命令行参数 得到 配置文件路径
)

// 初始化线程数量，取当前电脑核数
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 初始化命令行读取功能
func initArgs() {

	//1,保存的地址,2配置命令的名字,3默认值,4,使用介绍
	//eg：master -config ./master.json

	flag.StringVar(&confFile, "config",
		"./master.json", "制定配置文件：master.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)

	// 初始化命令行参数 4
	initArgs()

	// 初始化线程 1
	initEnv()

	// 加载配置 3
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}
	//  启动 任务管理器 5，要在 HTTP服务 前
	if err = master.InitJobMgr(); err != nil {
		fmt.Println(err) // 这里打印 可以确定错误位置
		goto ERR
	}

	// 启动API HTTP服务 2
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 主协程 不能退出 6
	for {
		time.Sleep(1 * time.Second)
	}

	return
ERR:
	fmt.Println(err)
}
