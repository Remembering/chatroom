package process

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gocode/project/chatroom2/client/global"
	"gocode/project/chatroom2/client/utils"
	"gocode/project/chatroom2/common/message"
)

type SmsProcess struct{}

// 发送群聊的消息
func (smsProc *SmsProcess) SendGroupMes(content string) (err error) {

	//1. 创建一个Mes
	var Mes message.Message
	Mes.Type = message.SmsMesType

	//2. 创建一个SmsMes 实例
	var smsMes message.SmsMes
	smsMes.Content = content //内容
	smsMes.UserMobile = CurUser.UserMobile
	smsMes.UserStatus = CurUser.UserStatus

	//3. 序列化Mes
	date, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal() err=", err.Error())
		return
	}

	Mes.Date = string(date)

	//4. 对Mes再次序列化
	date, err = json.Marshal(Mes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal() err=", err.Error())
	}

	//5. 将序列化的Mes发送给服务器

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6. 发送
	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("SendGroupMes err=", err.Error())
	}
	return
}

// 输出网页服务器里所有用户的信息
func (smsProc *SmsProcess) GetUserList(userMobile string) {
	//实例一个request对象并配置请求方法,

	req, _ := http.NewRequest("GET", "http://localhost:8021/u/v1/user/list", nil)
	//Header添加Content-Type
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-token", global.MobileAndToken[userMobile])
	//使用默认client发送请求
	rsp, _ := http.DefaultClient.Do(req)
	//读取网页返回body的内容
	body, _ := io.ReadAll(rsp.Body)
	//创建用户信息的实例
	var users []struct {
		Id     int    `json:"id"`     //数据库的id
		Name   string `json:"name"`   //昵称
		Mobile string `json:"mobile"` //手机号
	}
	//把body的内容写入rspjson
	json.Unmarshal(body, &users)
	//打印用户信息
	fmt.Println("用户列表如下:")
	for _, user := range users {
		fmt.Println("id:", user.Id)
		fmt.Println("mobile:", user.Mobile)
		fmt.Println()
	}
}

// 发送下线消息
func (smsProc *SmsProcess) SendExitMes(mobile string) (err error) {
	//1. 创建一个Mes
	var Mes message.Message
	Mes.Type = message.ExitMesType

	//2. 创建一个SmsMes 实例
	var smsMes message.ExitMes
	smsMes.Mobile = mobile

	//3. 序列化Mes
	date, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal() err=", err.Error())
		return
	}

	Mes.Date = string(date)

	//4. 对Mes再次序列化
	date, err = json.Marshal(Mes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal() err=", err.Error())
	}

	//5. 将序列化的Mes发送给服务器

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6. 发送
	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("SendGroupMes err=", err.Error())
	}
	return
}
