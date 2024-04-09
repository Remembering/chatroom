package process

import (
	"encoding/json"
	"fmt"

	"gocode/project/chatroom2/common/message"
)

func outputGroupMes(mes *message.Message) {

	//显示即可
	//1. 反序列化mes.Date
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Date), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//显示信息
	info := fmt.Sprintf("用户%s:\t%s \n:%s", smsMes.UserName, smsMes.UserMobile, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
