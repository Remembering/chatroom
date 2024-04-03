package process

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gocode/project/chatroom2/client/global"
	"gocode/project/chatroom2/client/utils"
	"gocode/project/chatroom2/common/message"
)

// 关联一个用户登陆的struct
type UserProcess struct{}

// 获取数字验证码
// 返回值字符串
// 通过post请求获取
func GetChaptcha() (string, string) {
	//实例一个request对象并配置请求方法,
	req, _ := http.NewRequest("GET", "http://localhost:8021/u/v1/base/captcha", nil)
	//使用默认client发送请求
	rsp, _ := http.DefaultClient.Do(req)
	//状态码不等于200请求失败的逻辑
	if rsp.StatusCode != http.StatusOK {
		fmt.Println("[GetChaptcha] 获取验证码失败")
		return "", ""
	}
	//读取网页返回body的内容
	body, _ := io.ReadAll(rsp.Body)
	//实例map对象接收body内容, interface{}可以接受任意类型的内容
	rspJson := map[string]interface{}{}
	//把body的内容写入rspjson
	json.Unmarshal(body, &rspJson)
	//返回string 并对interface{}断言为string
	return rspJson["answer"].(string), rspJson["id"].(string)
}

// 发送手机短信,通过网页请求获取
func PostSms(userMobile string) error {
	//发送信息的实例
	var sendSmsInfo struct {
		Mobile string
		Type   uint
	}
	//赋值手机号
	sendSmsInfo.Mobile = userMobile
	//赋值type
	sendSmsInfo.Type = uint(1)
	//将实例进行json序列化
	data, _ := json.Marshal(&sendSmsInfo)
	//实例一个request对象并配置请求方法,
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8021/u/v1/base/send_sms", bytes.NewBuffer(data))
	//Header添加Content-Type
	req.Header.Add("Content-Type", "application/json")
	//使用默认client发送请求
	rsp, _ := http.DefaultClient.Do(req)
	//读取网页返回body的内容
	body, _ := io.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	//返回状态码不是http.StatusOK就错误
	if rsp.StatusCode != http.StatusOK {
		fmt.Println("[PostSms] 错误" + string(body))
		return errors.New("[PostSms] 错误")
	}
	return nil
}

// 注册界面
// 向服务器发送注册信息
func (userProc *UserProcess) Register(userMobile string, userPwd string) (err error) {
	//1. 获取对服务器的client
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2. 发送手机短信进行注册
	// 输入y或Y确认发送短信, 否则退出注册
	var k string
	fmt.Println("输入y或Y发送手机短信")
	fmt.Scanf("%s", &k)
	//根据提示输入y或者Y就进行发送短信逻辑
	if k == "y" || k == "Y" {
		if err = PostSms(userMobile); err != nil {
			return
		}
	} else {
		return
	}
	code := ""
	fmt.Println("请输入手机短信的验证码(按回车继续):")
	fmt.Scanf("%s", &code)

	//2. 发送网页注册请求
	//4.创建一个RegisterMes 结构体 并赋值
	var registerMes message.RegisterMes
	registerMes.UserMobile = userMobile
	registerMes.UserPwd = userPwd
	registerMes.Code = code
	//将passWordLoginForm json化
	data, _ := json.Marshal(&registerMes)
	//实例一个req的Request对象并配置好请求格式和发送的内容
	req, _ := http.NewRequest("POST", "http://localhost:8021/u/v1/user/register", bytes.NewBuffer(data))
	//请求使用 application/json
	req.Header.Add("Content-Type", "application/json")
	//发送Post请求进行登陆
	res, _ := http.DefaultClient.Do(req)
	//读取body的内容
	body, _ := io.ReadAll(res.Body)
	//请求状态不是OK 就登陆失败的逻辑
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("[Register] POST请求 失败")
		fmt.Println(string(body))
		return
	}

	//3.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType //发送消息的类型

	//4.将registerMes 序列化
	date, err := json.Marshal(registerMes) //放回date类型为[]byte
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//5.把date赋给 mes.Date字段
	mes.Date = string(date) //获得RegisterMes结构体的值的字符串类型

	//6.将mes进行序列化
	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//11 这里还需要处理服务器端放回的消息
	//创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//发送date给服务器端
	err = tf.WritePkg(date)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
		return
	}

	//阻塞并等待接收服务器传来的信息
	mes, err = tf.ReadPkg() // mes 就是 RegisterResMes
	if err != nil {
		fmt.Println("ReadPkg() err=", err)
		return
	}

	//将mes的date部分反序列化成 RegisterResMes实例
	var RegisterResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Date), &RegisterResMes)
	//返回状态码200为成功
	if RegisterResMes.Code == 200 {
		fmt.Println("注册成功了,你重新登录一把")
		os.Exit(0)
	} else {
		//否则报错
		fmt.Println(RegisterResMes.Error)
		os.Exit(0)
	}
	return
}

// 登陆界面
// 向服务器发送登陆信息
func (userProc *UserProcess) Login(userMobile string, userPwd string) (err error) {

	//1.链接到服务器端
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("netDial err", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2. 获取验证码
	chaptcha, chaptchaId := GetChaptcha()
	//返回为空字符串,有错误
	if chaptcha == "" {
		return
	}
	//提示用户输入验证码
	fmt.Printf("请输入验证码:\t%s\n", chaptcha)
	//接收用户的按键输入的实例
	reChaptcha := ""
	fmt.Scanf("%s", &reChaptcha)
	//输入的和系统给的不一样则错误
	if chaptcha != reChaptcha {
		fmt.Println("输入的验证码错误")
		return
	}

	//3. 发送请求
	//创建发送登陆信息的实例并赋值
	loginMes := message.LoginMes{
		Mobile:    userMobile,
		PassWord:  userPwd,
		Captcha:   chaptcha,
		CaptchaId: chaptchaId,
	}
	//将passWordLoginForm json化
	data, _ := json.Marshal(&loginMes)
	//实例一个req的Request对象并配置好请求格式和发送的内容
	req, _ := http.NewRequest("POST", "http://localhost:8021/u/v1/user/pwd_login", bytes.NewBuffer(data))
	//请求使用 application/json
	req.Header.Add("Content-Type", "application/json")
	//发送Post请求进行登陆
	res, _ := http.DefaultClient.Do(req)
	//读取body的内容, 网页服务器返回的内容
	body, _ := io.ReadAll(res.Body)
	//请求状态不是OK 就登陆失败的逻辑
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("[Login] POST请求 失败")
		fmt.Println(string(body))
		return
	}
	//创建实例接收网页服务器返回的内容
	var role struct {
		Role  int32  `json:"role"`  //网页服务器返回该登陆用户的角色,普通用户或者是管理员
		Token string `json:"token"` //返回token与该用户绑定的信息,说明该用户正常
	}
	//将body的Json化内容反序列化赋值给role
	json.Unmarshal(body, &role)

	//4.准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType //发送消息的类型，与服务器的约定

	//6.将loginMes 序列化
	date, err := json.Marshal(loginMes) //放回date类型为[]byte
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//7.把date赋给 mes.Date字段
	mes.Date = string(date) //获得LoginMes结构体的值的字符串类型

	//8.将mes进行序列化
	date, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//9.到这个时候 ,data存储了mes的序列化，mes有发送信息的数据类型还有内容， date就是我们要发送到的消息
	//10.1先把date的长度发送给服务器
	//先获取到date的长度->转成一个表示长度的切片
	pkgLen := uint32(len(date))
	var buf [4]byte                              //为什么长度假设先4呢,uint32的大小
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) //把pkgLen的值转成[]byte类型的序列
	//相当于 buf[0:4] = []byte(pkgLen)
	//发送长度
	n, err := conn.Write(buf[0:4]) //发的是buf存储pkgLen值的切片
	if err != nil || n != 4 {
		fmt.Println("conn.Write(buf) err=", err)
		return
	}

	_, err = conn.Write(date) //发送序列化mes
	if err != nil {
		fmt.Println("conn.Write(date) err=", err)
		return
	}

	//11 这里还需要处理服务器端放回的消息
	//创建一个Transfer 实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//阻塞并等待接收服务器传来的信息
	mes, err = tf.ReadPkg() //
	if err != nil {
		fmt.Println("ReadPkg() err=", err)
		return
	}
	//创建聊天服务器返回信息的数据类型实例
	var loginResMes message.LoginResMes
	//将mes的date部分反序列化成 LoginResMes
	err = json.Unmarshal([]byte(mes.Date), &loginResMes)
	//返回状态码200为成功
	if loginResMes.Code == 200 {
		//保存登陆的token
		global.MobileAndToken[userMobile] = role.Token

		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserMobile = userMobile
		CurUser.UserStatus = message.UserOnline

		//现在可以显示当前在线用户列表
		fmt.Println("当前在线用户列表如下:")
		//遍历服务器返回的在线信息用户的数组
		for _, mobile := range loginResMes.UsersMobile {
			//创建所有在线用户的实例并赋值
			user := &message.User{
				UserMobile: mobile,
				UserStatus: message.UserOnline,
			}
			// 保存 在 客户端的 onlineUsers
			onlineUsers[mobile] = user

			//如果我们要求不显示自己在线， 下面我们增加一段代码
			if mobile == userMobile {
				continue
			}
			//打印用户信息
			fmt.Println("用户mobile:\t", mobile)
		}

		//3次换行效果
		fmt.Println()
		fmt.Println()
		fmt.Println()

		//这里我们还需要在客户端启动一个协程(类似子进程, 但占用资源小,轻量型的子进程)
		//该协程保持和服务器端的一个通讯..如果服务器有数据推送给客户端
		//则接送并写显示在客户端的终端。
		go serverProcessMes(conn, userMobile)

		//异常信号处理
		go func() {
			//初始化通道类型的实例并创建一个该类型大小的内存空间
			quit := make(chan os.Signal, 1)
			//遇到按键control+c等退出的singal就继续执行逻辑
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			//阻塞流程控制, 直到触发信号接触阻塞
			<-quit
			//提示退出
			fmt.Println("你选择退出系统!...")
			//退出时向服务器发送下线通知,服务器维护在线用户的map中存在该用户,就删掉退出用户
			smsProcess := &SmsProcess{}
			smsProcess.SendExitMes(userMobile)
			//退出客户端
			os.Exit(0)
		}()

		//1.显示登陆成功后的菜单.. 2:管理员界面	 1:普通用户界面
		for {
			//管理员界面
			if role.Role == 2 {
				ShowMenu2(userMobile)
			}
			//普通用户界面
			if role.Role == 1 {
				ShowMenu(userMobile)
			}

		}

	} else {
		//输出错误
		fmt.Println(loginResMes.Error)
	}
	return
}
