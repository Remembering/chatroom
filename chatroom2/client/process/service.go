package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"gocode/project/chatroom2/client/utils"
	"gocode/project/chatroom2/common/message"
)

// 显示登陆成功后普通用户的界面......
func ShowMenu(userMobile string) {
	fmt.Printf("-------------恭喜%s登陆成功------------\n", userMobile)
	fmt.Println("-----------1. 显示在线用户列表-----------")
	fmt.Println("-----------2. 发送消息------------------")
	fmt.Println("-----------3. 信息列表------------------")
	fmt.Println("-----------4. 退出系统------------------")
	fmt.Println("---------------请选择(1-4)--------------")

	//接收用户的选择
	var key int

	//接收用户聊天的按键输入
	var content string

	//因为， 我们总会使用到smsProcess实例， 因此我们将其定义在switch外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	//根据输入进入不同界面
	switch key {
	case 1:
		outputOnlineUser(userMobile)
	case 2:
		fmt.Println("你想对大家说的什么:)\t\t\t(输入exit退出聊天)")
		for content != "exit" {
			fmt.Scanln(&content)
			smsProcess.SendGroupMes(content)
		}
	case 3:
		showMesage(userMobile)
	case 4:
		fmt.Println("你选择退出系统!...")
		//退出时维护在线用户的map中存在该用户,就删掉退出用户
		smsProcess.SendExitMes(userMobile)
		os.Exit(0)
	default:
		fmt.Println("您输入的选项有误!...")
	}
}

// 显示登陆成功后管理员的界面......
func ShowMenu2(userMobile string) {
	fmt.Printf("-------------恭喜%s登陆成功------------\n", userMobile)
	fmt.Println("------------1. 显示用户列表--------------")
	fmt.Println("------------2. 退出--------------------")
	fmt.Println("-----------------请选择-----------------")
	var key int
	//因为， 我们总会使用到smsProcess实例， 因此我们将其定义在switch外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		smsProcess.GetUserList(userMobile)
	case 2:
		fmt.Println("你选择退出系统!...")
		//退出时维护在线用户的map中存在该用户,就删掉退出用户
		smsProcess.SendExitMes(userMobile)
		os.Exit(0)
	default:
		fmt.Println("您输入的选项有误!...")
	}
}

// 和服务器端保持通讯
// 参数需要服务器的连接标志,和当前用户的手机号
func serverProcessMes(conn net.Conn, mobile string) {
	//创建一个Transfer 实例， 不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//如果读取到了消息， 又是下一步处理逻辑
		//根据消息的类型进行不同的逻辑
		switch mes.Type {
		// 有人上线了
		case message.NotifyUserStatusMesType:
			//1. 取出.NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Date), &notifyUserStatusMes)
			//2. 把这个用户信息， 状态保存到客户端map[string]User中
			updateUserStatus(&notifyUserStatusMes, mobile)
		// 有人群发消息
		case message.SmsMesType:
			outputGroupMes(&mes)
		//有人下线了
		case message.ExitMesType:
			//1. 取出exitMes
			var exitMes message.ExitMes
			json.Unmarshal([]byte(mes.Date), &exitMes)
			//2. 把这个用户信息， 状态删除到客户端map[string]User中
			UserExitFunc(exitMes.Mobile)
		default:
			fmt.Println("服务器返回了未知的消息类型")
		}
	}
}

// 输出登陆用户的信息
func showMesage(userMobile string) {
	fmt.Println("信息列表")
	if user, ok := onlineUsers[userMobile]; ok {
		fmt.Println("昵称:", userMobile)
		fmt.Println("手机号:", user.UserMobile)
		fmt.Println("状态:", "在线")
	} else {
		fmt.Println("[信息列表] 错误")
	}
}
