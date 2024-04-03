package response

import (
	"fmt"
	"time"
)

// 自动移到time类型, 为了符合json内部转换到语法
type JsonTime time.Time

// 字段名转换json的键时会自动调用
// 将当前时间格式化为2006-01-02的样式
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	BirthDay JsonTime `json:"birthday"`
	// Birthday string `json:"birthday"`
	Gender string `json:"gender"`
	Mobile string `json:"mobile"`
}
