/*
time:2024.2.25
聊天室的主菜单显示
*/
package main

import (
	"fmt"

	"gocode/project/chatroom2/client/process"
)

// 定义两个变量，一个用户ID， 一个用户密码
var userMobile string
var userPwd string

// 主菜单函数
// 根据用户按键的输入进入不同的界面,登陆和注册
func main() {

	//接收用户的选择
	var key int

	//判断是否继续显示菜单
	loop := true

	//循环显示菜单
	for {
		fmt.Println("--------------------欢迎登陆多人聊天系统--------------------")
		fmt.Println("                      1.登陆聊天系统")
		fmt.Println("                      2.注册用户")
		fmt.Println("                      3.退出系统")
		fmt.Println("                      请选择(1-3)")

		//按键输入
		fmt.Scanf("%d", &key)

		//根据输入进入不同界面
		switch key {

		case 1: //登陆界面
			fmt.Println("登陆聊天室")
			fmt.Println("请输入账户Mobile")
			fmt.Scanf("%s\n", &userMobile)
			fmt.Println("请输入账户的密码")
			fmt.Scanf("%s\n", &userPwd)
			//1. 创建一个UserProcess实例
			up := &process.UserProcess{}
			//实例里的结构体绑定函数Login
			up.Login(userMobile, userPwd)

		case 2: //注册界面
			fmt.Println("注册用户")
			fmt.Println("请输入用户的Mobile:")
			fmt.Scanf("%s\n", &userMobile)
			fmt.Println("请输入用户的密码:")
			fmt.Scanf("%s\n", &userPwd)
			//2.调用UserProcess实例 完成注册的请求
			up := &process.UserProcess{}
			//实例里的结构体绑定函数Register
			up.Register(userMobile, userPwd)

		case 3:
			fmt.Println("退出系统")
			//继续显示菜单的标志为否
			loop = false
		default:
			fmt.Println("您的输入有误!请重新输入!")
		}

		//每次显示菜单判断标志是否继续显示菜单
		if !loop {
			break
		}
	}

}
