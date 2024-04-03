package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	"gocode/project/chatroom2/common/message"
)

// 这里将这些方法关联到/封装到结构体中
type Transfer struct {
	// 分析它应该有哪些字段
	Conn net.Conn
	Buf  [8096]byte //传输时使用的缓冲
}

// 最终获得将对方发过来的信息保存在相应结构体中
func (tran *Transfer) ReadPkg() (mes message.Message, err error) {

	fmt.Println("读取客户端的数据。。。")
	//conn.Read() 在conn没有被关闭的情况下 ，才会阻塞
	//如果客户端关闭了 conn 则 ， 就不会阻塞
	_, err = tran.Conn.Read(tran.Buf[:4]) // 返回值int是读了多少字节(Len),然后把读的数据赋给buf[0:4]; 读不到会一直等
	if err != nil {
		return
	}

	//根据buf[:4] 转成一个uint32类型
	pkgLen := binary.BigEndian.Uint32(tran.Buf[0:4]) //把buf的值转换为uint32

	//根据pkgLen 读取消息内容

	n, err := tran.Conn.Read(tran.Buf[:pkgLen]) //先获取切片的长度让后根据长度读的数据赋值给切片
	if n != int(pkgLen) || err != nil {
		return
	}

	//把pkgLen 反序列化成 -> message.Message
	// 技术就是一层窗户纸 &mes
	err = json.Unmarshal(tran.Buf[:pkgLen], &mes) //返回值的时候已经定义了; 把buf里读取的内容反序列化后的值赋给mes指向的结构体
	if err != nil {
		fmt.Println("json.Unmarshal() err=", err)
		return
	}
	return
}

func (tran *Transfer) WritePkg(date []byte) (err error) {

	//先发送一个长度给对方
	pkgLen := uint32(len(date))
	// var buf [4]byte//为什么长度假设先4呢，因为装得下吧
	binary.BigEndian.PutUint32(tran.Buf[0:4], pkgLen) //把pkgLen的值转成[]byte类型的序列
	//相当于 buf[0:4] = []byte(pkgLen)
	//发送长度
	n, err := tran.Conn.Write(tran.Buf[0:4]) //发的是buf存储pkgLen值的切片
	if err != nil || n != 4 {
		fmt.Println("conn.Write(buf) err=", err)
		return
	}

	//发送date
	n, err = tran.Conn.Write(date) //发的是buf存储pkgLen值的切片
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Write(buf) err=", err)
		return
	}

	return

}