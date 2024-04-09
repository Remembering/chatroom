package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gocode/project/chatroom_srvs/user_srv/model"
)

func main() {
	dsn := "root:00000000@tcp(127.0.0.1:3306)/chatroom_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&model.User{})
	// fmt.Println(genMd5("xxxxx_123456"))
	//将用户的密码变一下 随机字符串+用户密码
	//暴力破解 123456 111111 000000 彩虹表 盐值

	// // Using custom options

	//test创建用户
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(len(newPassword))
	fmt.Println(newPassword)

	passwordInfo := strings.Split(newPassword, "$")
	check := password.Verify("admin123", passwordInfo[2], passwordInfo[3], options)
	fmt.Println(check) // true
	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("mobb%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			Password: newPassword,
		}
		db.Save(&user)
	}
}

// func genMd5(code string) string {
// 	Md5 := md5.New()
// 	_, _ = io.WriteString(Md5, code)
// 	return hex.EncodeToString(Md5.Sum(nil))
// }
