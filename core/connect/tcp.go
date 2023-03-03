package connect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/io-game-go/common"
	"log"
	"net"
	"strconv"
)

type TcpSocket struct {
	MySocket
}

func (socket *TcpSocket) ListenBack(f func([]byte) []byte) {
	socket.proxyMethod = f
}

func (socket *TcpSocket) ListenAddr(addr string) {
	host, port, _ := net.SplitHostPort(addr)
	ip := net.ParseIP(host)
	parseInt, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		panic(err)
	}
	server := &net.TCPAddr{IP: ip, Port: int(parseInt)} //包含IP和端口
	listener, err := net.ListenTCP("tcp", server)
	if err != nil {
		panic(err)
	}
	// 延迟关闭监听器
	defer listener.Close()
	// 循环接收客户端的请求
	for {
		// 使用net.TCPListener的AcceptTCP方法，接收客户端的连接，返回一个net.TCPConn类型的对象和错误
		conn, err := listener.AcceptTCP()
		// 如果错误不为空，打印错误并继续循环
		if err != nil {
			log.Println("tcp连接错误:" + err.Error())
			continue
		}
		// 打印客户端的地址
		fmt.Println("client connected from:", conn.RemoteAddr())
		// 创建一个goroutine，调用处理连接的函数，传入连接对象作为参数
		go socket.handleConn(conn)
	}
}

func (socket *TcpSocket) handleConn(conn *net.TCPConn) {
	socket.AsyncResult(func(bytes []byte) {
		_, err := conn.Write(bytes)
		if err != nil {
			data := &MyProtocol{}
			data.SetData(bytes)
			_, err = conn.Write(Encode(data))
			if err != nil {
				common.FileLogger().Println("tcp写入错误:", err.Error())
				_ = conn.Close()
			}
		}
	})
	// 创建一个线程池，指定工作协程数为3，任务队列大小为10
	socket.pool = common.NewPool(20, 30)
	socket.pool.Start()

	// 延迟关闭连接
	defer conn.Close()
	// 创建一个缓冲区，用于存储接收到的数据
	buf := make([]byte, 4096)
	// 循环读取数据
	for {
		// 从连接中读取数据，返回读取的字节数和错误
		n, err := conn.Read(buf)
		if err != nil {
			common.FileLogger().Println("tcp读取错误: " + err.Error())
			break
		}
		// 调用解码函数，将字节切片转换为自定义协议的结构体
		mp := Decode(buf[:n])
		socket.InvokeMethod(mp.Data)

		//result := socket.proxyMethod(mp.Data)
		//if len(result) == 0 {
		//	continue
		//}
		//data := &MyProtocol{}
		//data.SetData(result)
		//_, err = conn.Write(Encode(data))
		//if err != nil {
		//	common.FileLogger().Println("tcp写入错误:", err.Error())
		//	break
		//}
	}
}

// MyProtocol 定义一个自定义协议的结构体，包含消息的长度、类型和内容
type MyProtocol struct {
	Length int32  // 消息的长度，用4个字节表示
	Data   []byte // 消息的内容，用字节切片表示，长度由Length决定
}

func (p *MyProtocol) SetData(data []byte) *MyProtocol {
	p.Length = int32(len(data))
	p.Data = data
	return p
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
