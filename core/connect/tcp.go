package connect

import (
	"bytes"
	"encoding/binary"
)

type Tcp struct {
}

func (t Tcp) ListenBack(f func([]byte) []byte) {
	//TODO implement me
	panic("implement me")
}

func (t Tcp) ListenAddr(addr string) {
	//TODO implement me
	panic("implement me")
}

// MyProtocol 定义一个自定义协议的结构体，包含消息的长度、类型和内容
type MyProtocol struct {
	Length int32  // 消息的长度，用4个字节表示
	Data   []byte // 消息的内容，用字节切片表示，长度由Length决定
}

func (p *MyProtocol) SetData(data []byte) {
	p.Length = int32(len(data))
	p.Data = data
}

// Encode 定义一个编码函数，将自定义协议的结构体转换为字节切片，用于发送数据
func Encode(mp *MyProtocol) []byte {
	// 创建一个缓冲区，用于存储编码后的数据
	buf := bytes.NewBuffer([]byte{})
	// 使用encoding/binary包中的Write函数，按照大端字节序，将结构体中的字段写入缓冲区
	_ = binary.Write(buf, binary.BigEndian, mp.Length)
	_ = binary.Write(buf, binary.BigEndian, mp.Data)
	// 返回缓冲区中的字节切片
	return buf.Bytes()
}

// Decode 定义一个解码函数，将字节切片转换为自定义协议的结构体，用于接收数据
func Decode(data []byte) *MyProtocol {
	// 创建一个缓冲区，用于存储解码后的数据
	buf := bytes.NewBuffer(data)
	// 创建一个自定义协议的结构体，用于存储解码后的字段
	mp := &MyProtocol{}
	// 使用encoding/binary包中的Read函数，按照大端字节序，从缓冲区中读取字段到结构体中
	_ = binary.Read(buf, binary.BigEndian, &mp.Length)
	// 根据Length的值，创建一个字节切片，用于存储Data字段
	mp.Data = make([]byte, mp.Length)
	_ = binary.Read(buf, binary.BigEndian, &mp.Data)
	// 返回解码后的结构体
	return mp
}
