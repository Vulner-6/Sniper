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
	"time"
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

/*
// 加载指定目录的全部插件，多个插件同时扫描多个url
func LoadAllPlugins(urlList []string, pocThread string, urlThread string) {
	// 类型转换
	pocThreadInt, err := strconv.Atoi(pocThread)

	if err != nil {
		fmt.Println("pocThread string 转换成 pocThreadInt int 失败")
	}

	urlThreadInt, err := strconv.Atoi(urlThread)
	if err != nil {
		fmt.Println("urlThread string 转换成 urlThreadInt int 失败")
	}

	url_num:=len(urlList);

	// 获取插件目录下的所有文件
	files, err := filepath.Glob(filepath.Join(pluginsDir, "*"+pluginSuffix()))
	file_num := len(files)

	if err != nil {
		fmt.Println("Error listing plugin files:", err)
		os.Exit(1)
	}
	if files == nil {
		fmt.Println("No plugins files ! Please check /plugins/compile/ directory !")
		os.Exit(1)
	}

	// 遍历加载并执行每个插件
	for _, file := range files {
		p, err := plugin.Open(file)
		if err != nil {
			fmt.Println("Error opening plugin:", err)
			continue
		}
		// 根据文件名，获取对应的标志
		file_name := getFileName(file) + ".go"
		poc_map := config.GetPocMapAndSoPath()
		symName := poc_map[file_name]
		// 根据插件定义的符号名获取插件对象
		symbol, err := p.Lookup(symName)

		if err != nil {
			fmt.Println("Error finding symbol:", err)
			continue
		}

		// 断言插件对象是否符合 Poc 接口
		poc, ok := symbol.(Poc)
		if !ok {
			fmt.Println("Plugin does not implement Poc interface")
			continue
		}

		// 执行插件
		result := poc.Run()
		fmt.Printf("Plugin %s result: %s\n", file, result)
	}
}
*/

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

// 从channel缓冲区中读取数据
func ReadData(ch chan string) {
	x := <-ch // 从Channel接收数据
	fmt.Println(x)
}

// 执行单个插件，扫描多个目标
func LoadOnePluginScanMore(targets []string, plugin_name string, target_thread string) {
	// 目标线程数字符串类型转整数类型

	targetThreadInt, err := strconv.Atoi(target_thread)
	if err != nil {
		fmt.Println("target_thread string 转换成 targetThreadInt int 失败！")
	}

	// 创建等待组，用于等待所有goroutine执行完毕
	var wg sync.WaitGroup
	// 创建用于收集结果的通道
	resultChan := make(chan string, targetThreadInt)
	// 创建计数器
	// thread_index := 0

	/*
		// 循环目标
		for i := 0; i < len(targets); i++ {
			// 计算缓冲区是否线程堆满
			if thread_index < targetThreadInt {
				wg.Add(1)
				go LoadOnePluginScanOne(targets[i], plugin_name, resultChan, &wg)
				thread_index = thread_index + 1
			} else {
				// 当缓冲区结果堆满后，必须将结果输出，才能继续填充缓冲区
				// 等待除最后一批外所有goroutine执行完毕,确保缓冲区有正确数量的数据
				wg.Wait()

				thread_index = 0
				// 打印结果通道中的结果
				for result := range resultChan {
					fmt.Println(result)
				}

			}

		}

	*/

	// 计算批次
	target_num := len(targets)
	/*
		cycle_remainder := target_num % targetThreadInt
		if cycle_remainder > 0 {
			cycle_num := target_num/targetThreadInt + 1
		} else {
			cycle_num := target_num / targetThreadInt
		}
		current_cycle_index := 0
	*/

	// 循环目标
	for i := 0; i < target_num; i++ {
		// 累计线程数
		wg.Add(1)
		// 执行子线程
		go LoadOnePluginScanOne(targets[i], plugin_name, resultChan, &wg)
		go ReadData(resultChan)
		/*
			// 处理通道缓冲区
			for r := 0; r < targetThreadInt; r++ {
				result := <-resultChan
				fmt.Println(result)
			}
		*/
	}

	// 等待所有goroutine执行完毕
	wg.Wait()

	time.Sleep(time.Second * 10)

	// 关闭结果通道，表示所有结果已经收集完毕
	close(resultChan)
	fmt.Println("所有线程执行结束！")

}

// 执行多个插件，扫描单个目标
func LoadMorePluginScanOne() {

}

// 执行多个插件，扫描多个目标
func LoadMorePluginScanMore() {

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
