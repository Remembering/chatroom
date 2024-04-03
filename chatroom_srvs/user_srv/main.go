package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"gocode/project/chatroom_srvs/user_srv/handler"
	"gocode/project/chatroom_srvs/user_srv/proto"
)

func main() {
	//获取终端命令的参数, go run main.go ip=127.0.0.1. port=50051
	//没有参数 就获取默认值
	Ip := flag.String("ip", "127.0.0.1", "ip地址")
	Port := flag.Int("port", 50051, "端口地址")
	flag.Parse()
	//打印
	fmt.Println("ip:" + *Ip)
	fmt.Println("port:", *Port)

	//创建新的server实例
	server := grpc.NewServer()
	//注册grpc服务的的代码
	proto.RegisterUserServer(server, &handler.UserServer{})
	//创建连接实例
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
	if err != nil {
		panic("failed to Listen" + err.Error())
	}
	//将lis实例给grpc来监听
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc" + err.Error())
	}
}
