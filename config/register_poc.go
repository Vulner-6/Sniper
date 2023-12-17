package config

// 编写完 poc 后，在这里注册 poc 标志，并返回 poc 文件名和 poc 标志的映射关系
func GetPocMapAndSoPath() (map[string]string, string) {
	// 插件编译后的目录绝对路径，根据自己电脑上的路径修改
	const pluginsDir = "/root/GitProject/Sniper/plugins/compile"

	poc_map := map[string]string{
		// 填写对应的插件代码名称和标志
		"myplugin.go": "MyPlugin",
		"do.go":       "MyDo",
		"say.go":      "Poc_CVE_2021",
	}

	// 返回插件映射关系与插件目录路径
	return poc_map, pluginsDir
}
