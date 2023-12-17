package main

import (
	"fmt"
	"time"
)

// MyPlugin 是插件中实现 Worker 接口的类型
type cve_2023 struct{}

// DoWork 是 Worker 接口的实现
func (p cve_2023) Run(target string) string {
	fmt.Printf(target + " Plugin Worker starting\n")
	// 插件的具体实现...
	time.Sleep(time.Second * 5)
	fmt.Printf(target + " Plugin Worker done\n")
	// 若漏洞不存在，则返回""字符串
	return "发现目标：" + target + "存在cve_2023漏洞！"
}

// exported ，用于根据 symbol 寻找到当前插件，当前应该寻找的值为 MyPlugin
var MyPlugin cve_2023
