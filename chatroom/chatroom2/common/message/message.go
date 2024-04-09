package message

// 确定一些消息类型
const (
	LoginMesType            = "LoginMes"            //登陆的发送的类型
	LoginResMesType         = "LoginResMes"         //登陆时服务器的返回的类型
	RegisterMesType         = "RegisterMesType"     //注册的发送类型
	RegisterMesResMesType   = "RegisterMesResMes"   //注册时服务器的返回的类型
	NotifyUserStatusMesType = "NotifyUserStatusMes" //服务器通知其他用户该用户上线的类型
	ExitMesType             = "ExitMesType"         //退出的发送类型
	SmsMesType              = "SmsMes"
)

const (
	UserOnline = iota
	Useroffline
	UserBusyStatus
)

type LoginMes struct {
	Mobile    string `form:"mobile" json:"mobile"`         //手机号码
	PassWord  string `form:"password" json:"password"`     //密码
	Captcha   string `form:"captcha" json:"captcha"`       //验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id"` //与用户绑定的验证码id
}

type Message struct {
	Type string `json:"type"` //消息的类型
	Date string `json:"date"` //消息的内容
}

type ExitMes struct {
	Mobile string `json:"mobile"` //手机号码
}

type LoginResMes struct {
	Code        int      `json:"code"`        //返回状态码 500, 表示用户为注册 200 表示登陆成功
	UsersMobile []string `json:"usersMobile"` //增加字段， 保存用户Mobile的一个切片
	Error       string   `json:"error"`       // 放回错误信息
	Role        int      `json:"role"`        //1代表管理员, 2代表普通用户

}

type RegisterMes struct {
	User        //继承User的字段
	Code string `json:"code"` //手机短信验证码
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码 400, 表示用户已经占用 200 表示注册成功
	Error string `json:"error"` // 放回错误信息
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserMobile string `json:"userMobile"` //用户mobile
	Status     int    `json:"status"`     //用户状态
}

type SmsMes struct {
	Content string //消息内容
	User           //匿名结构体 继承
}
