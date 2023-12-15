package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// pocPlugin 实现了 Poc 接口
type pocPlugin struct{}

// PocPlugin 是插件需要导出的符号,会在加载插件的时候查询这个符号名
var PocPlugin pocPlugin

// Run 是 Poc 接口的实现
func (p pocPlugin) Run() string {

	url := []string{"www.baidu.com", "www.bing.com", "www.google.com"}

	for _, v := range url {
		currentThreadIndex := strconv.Itoa(runtime.NumCPU())
		fmt.Println("say.go 线程号：" + currentThreadIndex + "正在扫描:" + v)
		time.Sleep(2 * time.Second)
	}
	return "say.go 扫描成功"
}
