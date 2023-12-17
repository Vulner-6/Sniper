package utils

import (
	"Sniper/config"
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strconv"
	"sync"
)

// Poc 接口定义了插件的基本要求
type Poc interface {
	Run(target string) string
}

// 初始化插件映射map和插件编译后的目录
var all_plugins, pluginsDir = config.GetPocMapAndSoPath()

// pluginSuffix 返回当前平台的插件文件后缀
func pluginSuffix() string {
	switch runtime.GOOS {
	case "windows":
		return ".dll"
	case "linux":
		return ".so"
	default:
		return ".so"
	}
}

// 获取路径字符中最后文件的文件名，无后缀
func getFileName(path string) string {
	// Use the filepath.Base function to get the base name of the path
	base := filepath.Base(path)

	// Remove the file extension to get the filename
	filename := base[:len(base)-len(filepath.Ext(base))]
	return filename
}

// 执行单个插件，扫描单个目标
func LoadOnePluginScanOne(target string, plugin_name string, resultChan chan<- string, wg *sync.WaitGroup) {
	// 获取插件目录下的指定编译后的插件对象
	poc_file := plugin_name + ".so"
	poc_file_path := filepath.Join(pluginsDir, poc_file)
	p, err := plugin.Open(poc_file_path)
	if err != nil {
		fmt.Println("Error opening plugin:", err)
		// return ""
	}

	// 检查插件标识
	source_file := plugin_name + ".go"
	current_symbol := all_plugins[source_file]
	symbol, err := p.Lookup(current_symbol)
	if err != nil {
		fmt.Println("Error finding symbol:", err)
		// return ""
	}

	// 断言插件对象是否符合 Poc 接口
	poc, ok := symbol.(Poc)
	if !ok {
		fmt.Println("Plugin does not implement Poc interface")
		// return ""
	}

	// 执行插件
	var result string = poc.Run(target)

	resultChan <- result

	// 等待所有goroutine执行完毕
	wg.Done()
	// return result
}

// 从 channel 缓冲区中读取数据
func ReadData(ch chan string, read_channel_wg *sync.WaitGroup, output string) {
	x := <-ch // 从Channel接收数据
	if x == "" {
		fmt.Println("暂时未发现漏洞，获得空字符串,正在扫描中，请等待程序执行完毕....")
		read_channel_wg.Done()
	} else {
		result := x + "\n\n"
		// 将数据写入指定文件,好像不需要设置锁
		WriteFileByLine(output, result)
		fmt.Println(x)
		read_channel_wg.Done()
	}
}

// 执行单个插件，扫描多个目标
func LoadOnePluginScanMore(targets []string, plugin_name string, target_thread string, plugin_wg *sync.WaitGroup, output string) {
	// 目标线程数字符串类型转整数类型
	targetThreadInt, err := strconv.Atoi(target_thread)
	if err != nil {
		fmt.Println("target_thread string 转换成 targetThreadInt int 失败！")
	}

	// 创建等待组，用于等待所有goroutine执行完毕
	var wg sync.WaitGroup
	var read_channel_wg sync.WaitGroup
	// 创建用于收集结果的通道
	resultChan := make(chan string, targetThreadInt)

	// 目标数量
	target_num := len(targets)

	// 若目标数量为0，则退出程序
	if target_num == 0 {
		fmt.Println("目标数量为0，未发现目标，退出程序！")
		os.Exit(1)
	}

	// 判断目标数量是大于等于批次，还是小于批次
	if target_num >= targetThreadInt {
		// 若目标数量大于等于每批次数量
		// 计算余数
		cycle_remainder := target_num % targetThreadInt
		// 不考虑余数的情况下计算批次
		cycle_num := target_num / targetThreadInt
		// 当前经过循环后的目标切片下标位置
		var current_cycle_target_index = 0
		for i := 0; i < cycle_num; i++ {
			// 每批线程数根据获得的数字决定
			for j := 0; j < targetThreadInt; j++ {
				// 累计线程数
				wg.Add(1)
				read_channel_wg.Add(1)
				// 计算下标
				current_cycle_target_index = targetThreadInt*i + j
				// 执行子线程
				go LoadOnePluginScanOne(targets[current_cycle_target_index], plugin_name, resultChan, &wg)
				go ReadData(resultChan, &read_channel_wg, output)

				// 重置下标
				current_cycle_target_index = 0
			}

			// 等待所有goroutine执行完毕
			wg.Wait()
			read_channel_wg.Wait()
		}

		// 如果有余数
		if cycle_remainder > 0 {
			// 余数对应的目标切片下标
			target_index := cycle_num * targetThreadInt
			// 根据线程数最后的余数循环1批
			for i := 0; i < cycle_remainder; i++ {
				// 累计线程数
				wg.Add(1)
				read_channel_wg.Add(1)

				// 执行子线程,注意这里因为有余数，所以对应目标的切片下标就不再是这里的i了，而是需要自己计算
				go LoadOnePluginScanOne(targets[target_index], plugin_name, resultChan, &wg)
				go ReadData(resultChan, &read_channel_wg, output)
				target_index++

			}
			// 等待所有goroutine执行完毕
			wg.Wait()
			read_channel_wg.Wait()
		}

	} else {
		// 小于批次（目标数量小于设置的目标线程数）
		for i := 0; i < target_num; i++ {
			// 累计线程数
			wg.Add(1)
			read_channel_wg.Add(1)

			// 执行子线程
			go LoadOnePluginScanOne(targets[i], plugin_name, resultChan, &wg)
			go ReadData(resultChan, &read_channel_wg, output)
		}
		// 等待所有goroutine执行完毕
		wg.Wait()
		read_channel_wg.Wait()
	}

	// 关闭结果通道，表示所有结果已经收集完毕
	close(resultChan)
	plugin_wg.Done()
	fmt.Println("所有线程执行结束！")

}

// 执行多个插件，扫描多个目标
func LoadMorePluginScanMore(targets []string, plugins []string, target_thread string, plugin_thread string, output string) {
	// 字符串转整数类型
	pluginThreadInt, err := strconv.Atoi(plugin_thread)
	if err != nil {
		fmt.Println("plugin_thread string 转换成 pluginThreadInt int 失败！")
	}

	// 计算plugin数量
	plugin_num := len(plugins)

	// 声明插件多线程等待组
	var plugin_wg sync.WaitGroup

	// 若插件数量为0，则退出程序
	if plugin_num == 0 {
		fmt.Println("插件数量为0，未发现插件，退出程序！")
		os.Exit(1)
	}

	// 判断插件数量大于等于插件线程数量，还是小于等于插件线程数量
	if plugin_num >= pluginThreadInt {
		// 计算余数
		cycle_remainder := plugin_num % pluginThreadInt
		// 不考虑余数的情况下计算批次
		cycle_num := plugin_num / pluginThreadInt
		// 当前经过循环后的插件切片下标位置
		var current_cycle_plugin_index = 0

		for i := 0; i < cycle_num; i++ {
			// 每批线程数根据获得的数字决定
			for j := 0; j < pluginThreadInt; j++ {
				// 累计线程数
				plugin_wg.Add(1)

				// 计算下标
				current_cycle_plugin_index = pluginThreadInt*i + j
				// 执行子线程
				go LoadOnePluginScanMore(targets, plugins[current_cycle_plugin_index], target_thread, &plugin_wg, output)

				// 重置下标
				current_cycle_plugin_index = 0
			}

			// 等待所有goroutine执行完毕
			plugin_wg.Wait()

		}

		// 如果有余数
		if cycle_remainder > 0 {
			// 余数对应的插件切片下标
			plugin_index := cycle_num * pluginThreadInt
			// 根据线程数最后的余数循环1批
			for i := 0; i < cycle_remainder; i++ {
				// 累计线程数
				plugin_wg.Add(1)

				// 执行子线程,注意这里因为有余数，所以对应目标的切片下标就不再是这里的i了，而是需要自己计算
				go LoadOnePluginScanMore(targets, plugins[plugin_index], target_thread, &plugin_wg, output)

				plugin_index++

			}
			// 等待所有goroutine执行完毕
			plugin_wg.Wait()

		}

	} else {
		// 若插件数量小于插件线程数量
		for i := 0; i < plugin_num; i++ {
			// 累计线程数
			plugin_wg.Add(1)

			// 执行子线程
			go LoadOnePluginScanMore(targets, plugins[i], target_thread, &plugin_wg, output)
		}
		// 等待所有goroutine执行完毕
		plugin_wg.Wait()
	}

	fmt.Println("所有插件执行结束！")

}

func logTest() {
	// 创建一个logrus日志实例
	logger := logrus.New()

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	logger.Info("日志功能正在启动，下面测试不同级别的日志输出信息：")
	logger.Debug("This is a debug message.")
	logger.Info("This is an info message.")
	logger.Warn("This is a warning message.")
	logger.Error("This is an error message.")
	logger.Info("日志功能正常！")
}

// 一行一行读取文件，返回字符串切片
func ReadFileByLine(file_path string) []string {
	// 打开文件
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建一个字符串切片来存储文件的每一行
	var lines []string

	// 创建一个Scanner来读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 将每一行添加到字符串切片中
		lines = append(lines, scanner.Text())
	}

	// 检查Scanner是否发生错误
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	/*
		// 打印每一行
		for _, line := range lines {
			fmt.Println(line)
		}
	*/

	return lines
}

// 一一行行的写入字符串至文件
func WriteFileByLine(file_path string, input string) {
	// 打开文件，如果不存在则创建，文件权限为 0666，追加写入模式
	file, err := os.OpenFile(file_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// 追加写入字符串到文件
	_, err = file.WriteString(input)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
	fmt.Println("Content appended successfully.")
}
