package global

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 全局变量db
var (
	DB *gorm.DB
)

// init函数启动时自动运行
func init() {
	//连接数据库的格式user:pwd@tcp(locahost:port)/tableName?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:00000000@tcp(127.0.0.1:3306)/chatroom_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//创建logger实例, 用来配置终端的输出样式
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	var err error
	//创建db连接mysql的实例, 以及加入logger的配置
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //自动生成的名字是否需要复数
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}
