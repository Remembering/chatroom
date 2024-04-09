package model

import (
	"time"

	"gorm.io/gorm"
)

// 表的基本结构
type BaseModel struct {
	ID        int32          `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"column:add_time"`    //创建的时间
	UpdatedAt time.Time      `goem:"column:update_time"` //更新的时间
	DeletedAt gorm.DeletedAt //过期的时间
	IsDelete  bool           //是否物理删除
}

/*
	1. 密文 2. 密文不可反解
		1. 对新加密
		2. 非对称加密
		3. md5 信息商要算法
		密码如果不可以反解，用户找回密码 直接提供修改的链接
*/

// 用户表结构
type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"not null;type:varchar(100)"`
	NickName string     `gorm:"type:varchar(20)"`
	Brithday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"default:male;column:gender;type:varchar(6) comment 'female为女性 male为男性'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '2表示普通用户, 1表示管理员'"`
}
