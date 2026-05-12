package chat

// QueryReq 表示聊天消息查询请求。
type QueryReq struct {
	MessageId int64 `json:"messageId"` // 消息 ID
}

// MessageResp 表示聊天消息响应。
type MessageResp struct {
	MessageId int64  `json:"messageId"` // 消息 ID
	FromUid   int64  `json:"fromUid"`   // 发送者用户 ID
	Content   string `json:"content"`   // 消息内容
}

// CreateReq 表示聊天消息创建请求。
type CreateReq struct {
	ToUid   int64  `json:"toUid"`   // 接收者用户 ID
	Content string `json:"content"` // 消息内容
}

// ListReq 表示聊天消息列表请求。
type ListReq struct {
	Limit int `json:"limit"` // 返回数量上限
}

// ListResp 表示聊天消息列表响应。
type ListResp struct {
	Messages []*MessageResp `json:"messages"` // 消息列表
}
