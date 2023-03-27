package messages

type Message interface {
	GetMerge() int64                  //获取路由
	GetBody() []byte                  //获取数据
	GetHeartbeat() bool               //心跳
	GetCode() int64                   //获取状态码,0为成功
	GetMessage() string               //消息
	GetBytesResult() []byte           //转换为字节数据
	SetBody(data interface{}) Message // 其内部适应了两种类型, []byte 和对应实现的类型
	GetHeaders() string               //用于扩展其他参数
	Bind(v interface{}) (err error)   //绑定到对象上
	SetMerge(merge int64)             //设置路由
	SetCode(code int64)               //设置状态码
	SetMessage(message string)        //错误消息
}
