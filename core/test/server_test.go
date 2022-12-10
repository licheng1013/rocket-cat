package main

import (
	"core/core"
	"core/decoder"
	"core/message"
	"core/router"
	"log"
	"testing"
	"time"
)

func TestJsonServer(t *testing.T) {
	// 默认的消息实现: DefaultMessage
	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) interface{} {
		//log.Println("收到消息: ", string(msg.([]byte)))
		return msg
	})

	server := core.NewGameServer()
	server.Run()
}

func TestProtoServer(t *testing.T) {

	// 默认的消息实现: DefaultMessage
	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) interface{} {
		info := message.Info{}
		message.UnmarshalInterface(msg, &info)
		//log.Println(info.String())
		return &info
	})

	server := core.NewGameServer()
	// 设置编码器
	server.SetDecoder(decoder.ProtoDecoder{})
	server.Run()
}

func TestCore(t *testing.T) {
	var v = time.Now().UnixMilli()
	var count int64

	router.AddFunc(router.GetMerge(0, 1), func(msg interface{}) interface{} {
		return msg
	})
	decoder.SetDecoder(decoder.ProtoDecoder{})

	for true {
		count++
		info := message.Info{Info: "Ok"}
		protoMessage := message.ProtoMessage{Merge: router.GetMerge(0, 1), Body: message.MarshalBytes(&info)}
		// 编码解码
		merge, body := decoder.GetDecoder().DecoderBytes(message.MarshalBytes(&protoMessage))
		_ = router.ExecuteFunc(merge, body)

		startTime := time.Now().UnixMilli()
		if startTime-v > 1000 {
			log.Println("执行次数: ", count)
			count = 0
			v = startTime
		}
	}

}
