package process2

import (
	"encoding/json"
	"fmt"
	"net"

	"gocode/project/chatroom2/common/message"
	"gocode/project/chatroom2/service/utils"
)

type SmsProcess struct {
	//暂时不需要字段.....
}

// 群发消息
func (smsProc *SmsProcess) SendGropMes(mes *message.Message) {

	//遍历服务端的onlineUsers map[int]*UserProcess
	//将消息转发取出

	//取出mes的内容 smsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Date), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	date, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	for mobile, up := range userMgr.onlineUsers {
		//这时我们要过滤掉自己， 就是不发给自己
		if mobile == smsMes.UserMobile {
			continue
		}
		smsProc.SendMesToEachOnlineUser(date, up.Conn)
	}
}

// 把消息发给每一个在线的用户
func (smsProc *SmsProcess) SendMesToEachOnlineUser(date []byte, conn net.Conn) {

	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(date)
	if err != nil {
		fmt.Println("转发消息失败! err=", err)
		return
	}

}
