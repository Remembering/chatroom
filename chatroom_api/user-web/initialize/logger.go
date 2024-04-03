package initialize

import "go.uber.org/zap"

// 初始化日志输出
func InitLogger() {
	Logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(Logger)
}
