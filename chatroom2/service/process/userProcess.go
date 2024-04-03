package process2

import (
	"encoding/json"
	"fmt"
	"net"

	"gocode/project/chatroom2/common/message"
	"gocode/project/chatroom2/service/model"
	"gocode/project/chatroom2/service/utils"
)

type UserProcess struct {
	//应该有哪些字段？   考虑将来绑定关联的方法需要什么东西
	Conn net.Conn
	//增加一个字段表示该Conn是哪个用户的
	UserMobile string
}

// 这里我们编写通知所有在线的用户的一个方法
// userId 要通知其他在线用户， 我上线
func (userProc *UserProcess) NotifyothersOnlineUsers(userMobile string) {

	//遍历 onlineUsers, 然后一个一个的发送 NotifyUsersStatausMes
	for mobile, up := range userMgr.onlineUsers {
		//过滤掉自己
		if mobile == userMobile {
			continue
		}
		//开始通知[单独的写一个方法]
		//up很关键保存了conn
		up.NotifyMeOnline(userMobile)
	}
}

func (userProc *UserProcess) NotifyMeOnline(userMobile string) {
	//组装我们的NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserMobile = userMobile
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	date, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal() err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes 赋值给mes.Date
	mes.Date = string(date)
	//对mes再次序列化
	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal() err", err)
	}

	//发送,创建一个Transfer实例，发送
	tf := &utils.Transfer{
		Conn: userProc.Conn,
	}

	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("tf.WritePkg() err=", err)
		return
	}
}

func (userProc *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//核心代码....
	//1. 先从mes 中取出 mes.Date ， 并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Date), &registerMes)
	if err != nil {
		fmt.Println("json.UnMashal() fial err=", err)
		return
	}
	var resMes message.Message
	resMes.Type = message.RegisterMesType
	var registerResMes message.RegisterResMes
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "注册发生未知错误"
		}
	} else {
		registerResMes.Code = 200
		fmt.Println("注册成功!")
	}

	date, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Mashal() fial err=", err)
		return
	}

	resMes.Date = string(date)

	//5.对resMes序列化， 进行发送
	date, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Mashal() fail err=", err)
		return
	}

	//6.发送date 我们将其封装到WritePkg函数
	//因为使用了分层模式(mvc)， 我们先创建Transfer实例， 然后读取
	tf := &utils.Transfer{
		Conn: userProc.Conn,
	}

	err = tf.WritePkg(date)
	return
}

// 编写一个函数serverProcessLogin函数，专门处理登陆请求
func (userProc *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码....
	//1. 先从mes 中取出 mes.Date ， 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Date), &loginMes)
	if err != nil {
		fmt.Println("json.UnMashal() fial err=", err)
		return
	}
	//1。先声明一个 resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2.在声明一个 LoginResMes, 并完成赋值
	var loginResMes message.LoginResMes

	// 我们需要到redis数据库去完成验证
	// 1.使用model.MyUserDao到redis去验证
	user, err := model.MyUserDao.Login(loginMes.Mobile, loginMes.PassWord)
	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}

	} else {
		loginResMes.Code = 200
		loginResMes.Role = 1
		//这里，因为用户登陆成功，我们就把该登陆成功的用户放到userMgr中
		//将登陆成功的userId赋值给userProc
		userProc.UserMobile = loginMes.Mobile
		userMgr.AddOnlineUser(userProc)
		//通知其他在线用户，我上线了
		userProc.NotifyothersOnlineUsers(loginMes.Mobile)

		//将当前在线用户的mobile放到loginResMes.UsersMobile
		//遍历 UserMgr.onlineUsers
		for mobile := range userMgr.onlineUsers {
			loginResMes.UsersMobile = append(loginResMes.UsersMobile, mobile)
		}
		fmt.Println("用户:", user.UserMobile, "登陆成功")
	}

	//3。将loginResMes序列化
	date, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Mashal() fial err=", err)
		return
	}

	//4.将date赋值给resMes
	resMes.Date = string(date)

	//5.对resMes序列化， 进行发送
	date, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Mashal() fail err=", err)
		return
	}
	//6.发送date 我们将其封装到WritePkg函数
	//因为使用了分层模式(mvc)， 我们先创建Transfer实例， 然后读取
	tf := &utils.Transfer{
		Conn: userProc.Conn,
	}
	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("tf.WritePkg() err=", err)
		return
	}
	return
}
