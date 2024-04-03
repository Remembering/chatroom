//管理用户状态变化

package process2

import (
	"encoding/json"
	"fmt"

	"gocode/project/chatroom2/common/message"
)

type UserStatusProc struct{}

// 通知其他用户其中有一个用户下线了
func (u *UserStatusProc) UserExitFunc(mes *message.Message) {
	//创建退出消息的实例
	var exitMes message.ExitMes
	//把客户端发送的mes数据JSON反序列化赋值给exitMes
	json.Unmarshal([]byte(mes.Date), &exitMes)
	//维护在线用户, 删除该用户
	userMgr.DelOnlineUser(exitMes.Mobile)

	//打印下线用户
	fmt.Println("用户:", exitMes.Mobile, "下线")
	//发送给每一个在线的用户通知该用户下线
	data, _ := json.Marshal(exitMes)
	mes.Date = string(data)
	data, _ = json.Marshal(mes)
	smsProcess := &SmsProcess{}
	for mobile, up := range userMgr.onlineUsers {
		//这时我们要过滤掉自己， 就是不发给自己
		if mobile == exitMes.Mobile {
			continue
		}
		smsProcess.SendMesToEachOnlineUser(data, up.Conn)
	}
}
