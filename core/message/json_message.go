package message

import "gitee.com/licheng1013/go-util/common"

// JsonMessage 必须实现 Message 接口
// Json处理则必须先转换为json才能继续处理其他东西
type JsonMessage struct {
	Merge     int64  `json:"merge,omitempty"`
	Body      []byte `json:"body,omitempty"`
	Heartbeat bool   `json:"heartbeat,omitempty"`
	Code      int64  `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (j *JsonMessage) GetBytesResult() []byte {
	return []byte(common.JsonUtil.JsonToStr(j))
}

func (j *JsonMessage) SetBody(bytes []byte) {
	j.Body = bytes
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
