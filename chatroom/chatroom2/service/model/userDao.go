package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"

	"gocode/project/chatroom2/common/message"
)

// 我们在服务器启动后， 就初始化一个UserDao实例
// 把它做成一个全局的变量， 在需要和redis操作时就直接使用即可
var (
	MyUserDao *UserDao
)

//定义一个UserDao 结构体
//完成对User 结构体的各种操作..

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式， 创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 思考一下在UserDao 应该提供哪些方法给我们
// 1.根据用户ID返回一个User实例 + err
func (u *UserDao) getUserByMobile(conn redis.Conn, mobile string) (user *User, err error) {
	//通过给定的id去ridis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", mobile))
	if err != nil {
		if err == redis.ErrNil { // 在users 哈希中，没有找到对应的 mobile
			err = ERROR_USER_NOTEXISTS //自定义错误常量
		}
		return
	}
	user = &User{}
	//这里我们需要把res反序列化成一个user实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal() err=", err)
		return
	}
	return
}

//完成登陆的校验 Login
//1.Login 完成对用户的验证
//2.如果用户的id和pwd都正确，则返回一个user实例
//3.如果用户id或者pwd有错误，则返回一个对应的错误信息

func (u *UserDao) Login(userMobile string, userPwd string) (user *User, err error) {

	//先从UserDao的链接池中取出一个链接
	conn := u.pool.Get()
	defer conn.Close()
	user, err = u.getUserByMobile(conn, userMobile)
	if err != nil {
		return
	}
	//这是证明用户是获取到的
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (u *UserDao) Register(user *message.User) (err error) {

	//先从UserDao的链接池中取出一个链接
	conn := u.pool.Get()
	defer conn.Close()
	_, err = u.getUserByMobile(conn, user.UserMobile)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	//这时说明id在redis里面还没有， 则可以完成注册
	date, err := json.Marshal(user)
	if err != nil {
		return
	}
	//入库
	_, err = conn.Do("HSet", "users", user.UserMobile, string(date))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
