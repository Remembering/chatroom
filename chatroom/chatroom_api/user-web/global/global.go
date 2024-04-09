package global

import (
	"context"

	ut "github.com/go-playground/universal-translator"

	"gocode/project/chatroom_api/user-web/config"
)

var (
	// validator的全局翻译器
	Trans ut.Translator
	//通过viper读取本地配置文件内容的struct
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	//全局使用的上下分
	Ctx context.Context = context.Background()
)
