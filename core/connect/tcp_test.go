package connect

import (
	"fmt"
	"log"
	"net"
	"testing"
)

func TestMyProtocol(t *testing.T) {
	m := MyProtocol{}
	m.SetData([]byte("HelloWorld"))
	encode := Encode(&m)
	result := Decode(encode)
	fmt.Println(string(result.Data))
}

const message = "HelloWorld"

func TestTcpServer(t *testing.T) {
	channel := make(chan int)
	socket := TcpSocket{}
	go func() {
		socket.ListenBack(func(bytes []byte) []byte {
			return bytes
		})
		socket.ListenAddr(Addr)
	}()
	go Client(channel)
	select {
	case ok := <-channel:
		log.Println(ok)
	}
}

func Client(channel chan int) {
	// 连接服务器
	conn, err := net.Dial("tcp", Addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go func() {
		for {
			m := &MyProtocol{}
			_, err := conn.Write(Encode(m.SetData([]byte(message)))) // 发送数据
			if err != nil {
				fmt.Println("写入错误:", err)
				break
			}
		}
	}()
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf[:]) // 接收数据
		if err != nil {
			fmt.Println("读取错误:", err)
			break
		}
		decode := Decode(buf[:n])
		fmt.Printf("获取数据: %s\n", string(decode.Data))
		channel <- 0
	}
}
