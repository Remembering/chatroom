package model

// 定义一个用户的结构体
type User struct {
	//确定字段信息
	//为了反序列化成功，我们必须保证
	//用户信息的json字符串key 和 结构体的字段对应一致， 因为是值传递
	UserMobile string `json:"mobile"`   //用户手机号
	UserPwd    string `json:"password"` //用户密码
	UserName   string `json:"userName"` //昵称
}
