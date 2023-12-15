/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Sniper/utils"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行插件",
	Long:  `用于运行插件前的前置命令`,
	// 调用 run 命令时执行
	Run: func(cmd *cobra.Command, args []string) {

		// 针对单个目标扫描
		if target != "" {
			// 若对单个url调用全部插件扫描
			if plugins == "*" {
				fmt.Println("准备调用" + plugins + "插件。。。")

			} else {

				// 判断插件数量，若只有一个插件
				poc_slice := strings.Split(plugins, ",")
				if len(poc_slice) > 1 {
					// 调用多个插件对一个目标进行扫描

				} else {
					// 若对单个目标调用指定1个插件扫描
					if len(poc_slice) > 0 {
						// utils.LoadOnePluginScanOne(target, poc_slice[0])
					} else {
						fmt.Println("请至少输入一个待扫描的插件名称！")
					}
				}

			}
		}

		// 针对多个目标进行扫描
		if file != "" {
			// 若对多个url调用全部插件扫描
			if plugins == "*" {

			} else {
				// 若对多个url调用指定插件扫描
				// 获取多个扫描目标
				all_target := utils.ReadFileByLine(file)
				utils.LoadOnePluginScanMore(all_target, plugins, target_thread)
			}
		}

		// 一个插件扫描多个目标

		// 多个插件扫描多个目标

	},
}

// 定义 run 命令后的参数
var plugins string
var target string
var file string
var target_thread string
var poc_thread string
var output string

func init() {
	rootCmd.AddCommand(runCmd)

	// 获取当前时间戳，制作结果文件名
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	timestampStr := fmt.Sprintf("%d", timestamp)
	result_file_name := timestampStr + "_" + "result.txt"

	// Here you will define your flags and configuration settings.
	runCmd.Flags().StringVarP(&plugins, "plugins", "", "", "指定要运行的插件文件名称，用英文逗号分割。如：test,say （省略文件后缀）,*表示加载全部插件")
	runCmd.Flags().StringVarP(&target, "target", "", "", "指定单个扫描目标，扫描目标格式：协议+主机+[端口]")
	runCmd.Flags().StringVarP(&file, "file", "", "", "读取txt文件,加载多个扫描目标，扫描目标格式：协议+主机+[端口]")
	runCmd.Flags().StringVarP(&target_thread, "target_thread", "", "16", "指定一个插件扫描批量目标时的并发线程数，默认一个插件同时扫描16个url。")
	runCmd.Flags().StringVarP(&poc_thread, "poc_thread", "", "16", "指定同时多个插件扫描批量目标时的并发线程数，默认16个插件同时扫描--url_thread指定的url数量。")
	runCmd.Flags().StringVarP(&output, "output", "", result_file_name, "指定输出路径，默认输出至当前路径。")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
