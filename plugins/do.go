package main

import (
	"fmt"
	"runtime"
	"time"
)

// pocPlugin 实现了 Poc 接口
type testdopoc struct{}

// 标志，用于加载插件时寻找的标志变量，首字母必须大写，否则找不到
var Isitdo testdopoc

func (p testdopoc) Run() string {
	for i := 1; i < 6; i++ {
		fmt.Println("hello：当前线程号：", runtime.NumCPU(), "\n")
		time.Sleep(3 * time.Second)
	}
	return "do.go 执行完毕！"
}
