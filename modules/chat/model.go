package chat

type QueryReq struct {
	MessageId int64 `json:"messageId"`
}

type MessageResp struct {
	MessageId int64  `json:"messageId"`
	FromUid   int64  `json:"fromUid"`
	Content   string `json:"content"`
}

type CreateReq struct {
	ToUid   int64  `json:"toUid"`
	Content string `json:"content"`
}

type ListReq struct {
	Limit int `json:"limit"`
}

type ListResp struct {
	Messages []*MessageResp `json:"messages"`
}
