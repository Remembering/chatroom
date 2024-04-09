package process

// 维护本地onluneUSers ,有其他用户下线,删除该用户
func UserExitFunc(mobile string) {
	delete(onlineUsers, mobile)
}
