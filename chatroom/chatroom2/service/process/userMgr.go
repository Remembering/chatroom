package process2

import (
	"fmt"
)

//1.因为UserMgr实例在服务器端有且只有一个
//2.因为在很多地方都会使用到
//所以我们将其定义为一个全局变量

var (
	userMgr *UserMgr //协程之间共享堆空间
)

type UserMgr struct {
	onlineUsers map[string]*UserProcess
}

// 完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[string]*UserProcess, 1024),
	}
}

// 完成对onlieUsers的添加 修改
func (u *UserMgr) AddOnlineUser(up *UserProcess) {
	u.onlineUsers[up.UserMobile] = up
}

// 删除
func (u *UserMgr) DelOnlineUser(userMobile string) {
	delete(u.onlineUsers, userMobile)
}

// 返回当前所有在线的用户
func (u *UserMgr) GetOnlineUser() map[string]*UserProcess {
	return u.onlineUsers
}

// 根据Id返回对应的值
func (u *UserMgr) GetOnlineUserById(userMobile string) (up *UserProcess, err error) {

	//如何从map中取出一个值，带检测的方式
	up, ok := u.onlineUsers[userMobile]
	if !ok { //说明你要查找的这个用户， 当前不在线
		err = fmt.Errorf("用户%s 不存在或者不在线", userMobile) //把一个格式化的err返回去 赋值
		return
	}
	return
}
