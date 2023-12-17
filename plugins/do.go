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

func (p testdopoc) Run(target string) string {
	for i := 1; i < 6; i++ {
		fmt.Println("hello：当前线程号：", runtime.NumCPU(), "\n")
		time.Sleep(3 * time.Second)
		return ""
	}
	fmt.Println(target + "未发现漏洞！")
	return ""
}

// exported ，用于根据 symbol 寻找到当前插件，当前应该寻找的值为 MyDo
var MyDo testdopoc
