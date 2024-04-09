package main

import (
	"fmt"
	"io"
	"net"

	"gocode/project/chatroom2/common/message"
	process2 "gocode/project/chatroom2/service/process"
	"gocode/project/chatroom2/service/utils"
)

// 先创建一个Processor 的结构体
type Processor struct {
	Conn net.Conn
}

// 编写一个ServiceProcessMes 函数
// 功能： 根据客户端发送的消息的种类不同决定调用哪个函数来处理
func (proc *Processor) ServiceProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登陆的逻辑
		//创建一个UserProcess 实例
		up := &process2.UserProcess{
			Conn: proc.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		// 创建一个UserProcess
		up := &process2.UserProcess{
			Conn: proc.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGropMes(mes)
	case message.ExitMesType:
		userStatusProc := &process2.UserStatusProc{}
		userStatusProc.UserExitFunc(mes)
	default:
		fmt.Println("消息类型不存在, 无法处理!")
	}
	return
}

// 读取第一层解密的数据
func (proc *Processor) Process2() (err error) {

	//循环读客户端发送的信息
	for {
		//这里我们将读取数据包， 直接封装成一个函数ReadPkg(), 放回Message, Err
		//创建一个Transfer 实例 完成读包任务
		tf := &utils.Transfer{
			Conn: proc.Conn,
		}
		//阻塞流程控制,直到读取客户端的信息
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出")
				return err
			} else {
				fmt.Println("ReadPkg err", err)
				return err
			}
		}

		err = proc.ServiceProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
