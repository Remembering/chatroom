package main

import (
	"fmt"
	"net"
	"time"

	"gocode/project/chatroom2/service/model"
)

func process(conn net.Conn) {
	//这里要延时关闭conn
	defer conn.Close()

	//调用总控, 创建一个总控
	processor := &Processor{
		Conn: conn,
	}

	err := processor.Process2()
	if err != nil {
		fmt.Println("客户端和服务器端通讯协程错误 err=", err)
		return
	}
}

// 初始化redis数据库,以及redis的pool,init代表启动时自动运行
func init() {
	//这里我们编写一个函数， 完成对UserDao的初始化任务
	initPool("127.0.0.1:6379", 16, 0, 300*time.Second)
	initUserDao()
}

func initUserDao() {
	//这里的pool本身就是一个全局的变量  ...(redis)
	//这里需要注意初始化的顺序问题
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {

	//提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()

	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器。。。。。")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err", err)
		}
		//一旦链接成功，则启动一个协程和客户端保持通讯..
		go process(conn)
	}
}
