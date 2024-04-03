package process

import (
	"fmt"

	"gocode/project/chatroom2/client/model"
	"gocode/project/chatroom2/common/message"
)

// 客户端需要维护的map
var onlineUsers map[string]*message.User = make(map[string]*message.User, 10)
var CurUser model.CurUser // 我们在用户登陆成功后 完成对CurUser初始化,当前用户的临时变量

// 在客户端显示当前在线用户
func outputOnlineUser(userMobile string) {
	fmt.Println("当前在线用户列表:")
	for _, user := range onlineUsers {
		//过滤掉自己 号码与用户的号码相同就跳过这次循环
		if user.UserMobile == userMobile {
			continue
		}
		//打印信息
		fmt.Println("mobile:\t", user.UserMobile)
	}
}

// 编写一个方法,处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes, mobile string) {

	//适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserMobile]
	if !ok {
		user = &message.User{
			UserMobile: notifyUserStatusMes.UserMobile,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserMobile] = user
	outputOnlineUser(mobile)
}
