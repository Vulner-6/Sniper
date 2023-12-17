/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"Sniper/config"
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

		// 针对单个目标或少量两三个目标扫描
		if target != "" {
			// 无论一个还是多个目标，都将其变成切片
			target_slice := strings.Split(target, ",")
			if plugins == "*" {
				// 加载指定目录全部插件进行扫描
				// 初始化插件映射map和插件编译后的目录
				var all_plugins, _ = config.GetPocMapAndSoPath()
				// 提取所有的key
				plugins_slice := make([]string, 0, len(all_plugins))
				for key := range all_plugins {
					temp := strings.Split(key, ".")
					plugins_slice = append(plugins_slice, temp[0])
				}

				// 开始扫描
				utils.LoadMorePluginScanMore(target_slice, plugins_slice, target_thread, plugin_thread, output)

			} else {
				// 根据手动指定的1个或几个插件进行扫描
				plugins_slice := strings.Split(plugins, ",")
				utils.LoadMorePluginScanMore(target_slice, plugins_slice, target_thread, plugin_thread, output)
			}

		}

		// 针对多个目标进行扫描
		if file != "" {
			// 若对多个url调用全部插件扫描
			if plugins == "*" {
				// 初始化插件映射map和插件编译后的目录
				var all_plugins, _ = config.GetPocMapAndSoPath()
				// 获取多个扫描目标
				all_target := utils.ReadFileByLine(file)
				// 提取所有的key
				plugins_slice := make([]string, 0, len(all_plugins))
				for key := range all_plugins {
					temp := strings.Split(key, ".")
					plugins_slice = append(plugins_slice, temp[0])
				}

				// 开始扫描
				utils.LoadMorePluginScanMore(all_target, plugins_slice, target_thread, plugin_thread, output)

			} else {
				// 若对多个url调用指定插件扫描
				// 获取多个扫描目标
				all_target := utils.ReadFileByLine(file)
				all_plugin := strings.Split(plugins, ",")
				utils.LoadMorePluginScanMore(all_target, all_plugin, target_thread, plugin_thread, output)
			}
		}

	},
}

// 定义 run 命令后的参数
var plugins string
var target string
var file string
var target_thread string
var plugin_thread string
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
	runCmd.Flags().StringVarP(&plugin_thread, "plugin_thread", "", "16", "指定同时多个插件扫描批量目标时的并发线程数，默认16个插件线程同时扫描--target_thread指定的url数量。")
	runCmd.Flags().StringVarP(&output, "output", "", result_file_name, "指定输出路径，默认输出至当前路径。")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
