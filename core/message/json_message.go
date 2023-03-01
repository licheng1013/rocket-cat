package message

import (
	"github.com/io-game-go/common"
)

// JsonMessage 必须实现 Message 接口
// Json处理则必须先转换为json才能继续处理其他东西
type JsonMessage struct {
	Merge     int64  `json:"merge,omitempty"`
	Body      []byte `json:"body,omitempty"`
	Heartbeat bool   `json:"heartbeat,omitempty"`
	Code      int64  `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Headers   string `json:"headers,omitempty"`
}

func (j *JsonMessage) Bind(v interface{}) (err error) {
	if err = common.AssertPtr(v, "不是指针类型,无法绑定到结构体上"); err != nil {
		return
	}
	err = MsgKit.BytesToStruct(j.GetBody(), v)
	return
}

func (j *JsonMessage) GetHeaders() string {
	return j.Headers
}

func (j *JsonMessage) GetBytesResult() []byte {
	return MsgKit.StructToBytes(j)
}

func (j *JsonMessage) SetBody(bytes []byte) Message {
	j.Body = bytes
	return j
}

func (j *JsonMessage) GetMerge() int64 {
	return j.Merge
}

func (j *JsonMessage) GetBody() []byte {
	return j.Body
}

func (j *JsonMessage) GetHeartbeat() bool {
	return j.Heartbeat
}

func (j *JsonMessage) GetCode() int64 {
	return j.Code
}

func (j *JsonMessage) GetMessage() string {
	return j.Message
}
