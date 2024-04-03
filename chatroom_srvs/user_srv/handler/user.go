package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"gocode/project/chatroom_srvs/user_srv/global"
	"gocode/project/chatroom_srvs/user_srv/model"
	"gocode/project/chatroom_srvs/user_srv/proto"
)

// 完成grpc服务的代码的struct
// 因为proto自动生成的接口 有mustEmbedUnimplementedUserServer()
// 该匿名struct实现了该方法
// 所以加上了 proto.UnimplementedUserServer
type UserServer struct {
	proto.UnimplementedUserServer
}

// 把User的类型转换成UserInfoResponse类型
func ModelToResponse(user model.User) *proto.UserInfoResponse {
	//在grpc的message中字段默认值，你不能随便赋值ni1进去，容易出错
	//这里要搞清，哪些字段是有默议值的
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	//如果birthday的参数不为nil就赋值
	if user.Brithday != nil {
		userInfoRsp.BirthDay = uint64(user.Brithday.Unix())
	}
	return &userInfoRsp
}

// 分页功能
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 初始化User数组实例
	var users []model.User
	//db调用查询函数,获取所有用户并赋值给users
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	//初始化返回的数据类型实例并赋值
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	//实现分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	for _, user := range users {
		userInfoReq := ModelToResponse(user)
		rsp.Date = append(rsp.Date, userInfoReq)
	}
	return rsp, nil
}

// 通过id查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 初始化User实例
	var user model.User
	//db调用查询函数,获取所有用户并赋值给users
	result := global.DB.First(&user, req.Id)
	//返回数据为0提示错误
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户信息不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)

	return userInfoRsp, nil
}

// 通过Mobile查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 初始化User实例
	var user model.User
	//db调用条件查询函数,获取所有用户并赋值给users
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//返回数据为0提示错误
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户信息不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToResponse(user)

	return userInfoRsp, nil
}

// 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreatUserInfo) (*proto.UserInfoResponse, error) {
	// 初始化User实例
	var user model.User
	//db调用查询函数,获取所有用户并赋值给users
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	//返回数据为1代表已有该用户存在
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	//把网页服务器端请求的数据赋值
	user.NickName = req.NickName
	user.Mobile = req.Mobile
	user.Password = req.PassWord

	//密码加密
	//配置md5加密的设置 盐值长度16 迭代次数100 密码哈希长度
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	//db调用创建函数
	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	userInfoRep := ModelToResponse(user)
	return userInfoRep, nil
}

// 个人中心更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	// 初始化User实例
	var user model.User
	//db调用条件查询函数
	result := global.DB.First(&user, req.Id)
	//返回数据为0提示错误
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户信息不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	//把网页服务器端请求的数据赋值
	birthDay := time.Unix(int64(req.BirthDate), 0)
	user.NickName = req.NickName
	user.Brithday = &birthDay
	user.Gender = req.Gender
	result = global.DB.Save(user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

// 检查密码是否正确
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PassWordCheckInfo) (*proto.CheckResponse, error) {
	//获取之前的加密设置
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	// 将字符串通过自己里面的"$"分割成字符数组
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	//取出盐值 密码的哈希key ,验证密码是否匹配
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
