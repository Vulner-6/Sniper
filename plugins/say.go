package main

import (
	"fmt"
	"time"
)

// pocPlugin 实现了 Poc 接口
type cve_2021 struct{}

// Run 是 Poc 接口的实现
func (p cve_2021) Run(target string) string {
	fmt.Println("正在扫描:" + target)
	time.Sleep(2 * time.Second)
	return "cve_2021 扫描成功，存在漏洞cve_2021"
}

// Poc_CVE_2021 是插件需要导出的符号,会在加载插件的时候查询这个符号名
var Poc_CVE_2021 cve_2021
