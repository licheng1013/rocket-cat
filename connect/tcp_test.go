package connect

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestMyProtocol(t *testing.T) {
	m := MyProtocol{}
	m.SetData([]byte("HelloWorld"))
	encode := Encode(&m)
	result := Decode(encode)
	fmt.Println(string(result.Data))
}

func TestTcpServer(t *testing.T) {
	channel := make(chan int)
	socket := TcpSocket{}
	go func() {
		socket.ListenBack(func(uuid uint32, bytes []byte) []byte {
			return bytes
		})
		socket.ListenAddr(Addr)
	}()
	time.Sleep(time.Second / 2)
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
	go func() {
		data := HelloMsg
		for i := 0; i < 8; i++ {
			data += data
		}

		for {
			m := &MyProtocol{}
			_, err := conn.Write(Encode(m.SetData([]byte(data)))) // 发送数据
			if err != nil {
				panic(err)
			}
		}
	}()
	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf) // 接收数据
		if err != nil {
			panic(err)
		}
		decode := Decode(buf[:n])
		fmt.Println("获取数据: " + string(decode.Data))
		channel <- 0
	}
}
