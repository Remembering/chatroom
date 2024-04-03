package model

import (
	"net"

	"gocode/project/chatroom2/common/message"
)

// 因为客户端，我们很多地方会使用到curUser， 我们将其作为一个全局
type CurUser struct {
	Conn net.Conn //维护连接
	message.User
}
