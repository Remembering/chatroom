package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"gocode/project/chatroom_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn
var ctx context.Context = context.Background()

func Init() {
	var err error
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err = grpc.Dial(":50051", opt)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(ctx, &proto.PageInfo{
		Pn:    2,
		PSize: 2,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Date {
		fmt.Println(user.Id)
		checkRsp, err := userClient.CheckPassword(ctx, &proto.PassWordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err.Error() + "密码校验错误")
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestGetUserByMobile() {
	userRsp, err := userClient.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: "18782222224",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(userRsp.Mobile, userRsp.NickName, userRsp.PassWord)
}

func TestGetUserById() {
	userRsp, err := userClient.GetUserById(ctx, &proto.IdRequest{
		Id: 14,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(userRsp.Mobile, userRsp.NickName, userRsp.PassWord)
}

func TestUpdateUser() {
	_, err := userClient.UpdateUser(ctx, &proto.UpdateUserInfo{
		Id:        16,
		BirthDate: uint64(time.Now().Unix()),
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	TestGetUserList()
	// TestGetUserByMobile()
	// TestGetUserById()
	// TestUpdateUser()
	defer conn.Close()
}
