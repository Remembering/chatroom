package initialize

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"gocode/project/chatroom_api/user-web/global"
)

// 获取系统变量env
func GenEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 获取执行文件时所在的路径
func GetCurrentAbPathByExecutable() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = path[:index]
	return path
}

func InitConfig() {
	//存在该系统变量则返回具体的值, 否则返回false
	debug := GenEnvInfo("MSHOP_DEBUG")

	//共同的文件名字前缀
	configFilePrefix := "config"

	// path := GetCurrentAbPathByExecutable()
	// fmt.Println(path)
	// configFileName := fmt.Sprintf("%s/%s-pro.yaml", path, configFilePrefix)
	// if !debug {
	// 	configFileName = fmt.Sprintf("%s/%s-debug.yaml", path, configFilePrefix)
	// }
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}
	// New一个viper实例
	v := viper.New()

	//设置读取文件相对main.go文件的位置
	v.SetConfigFile(configFileName)

	//读取文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//将文件的内容解码到变量里
	v.Unmarshal(global.ServerConfig)

	//打印读取的信息
	zap.S().Infof("配置信息: %v", global.ServerConfig)

	// viper的功能 - 动态监控变化
	go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			zap.S().Infof("配置文件产生变化: %s", e.Name)
			v.ReadInConfig()
			v.Unmarshal(global.ServerConfig)
			zap.S().Infof("配置信息: %v", global.ServerConfig)
		})
	}()
}
