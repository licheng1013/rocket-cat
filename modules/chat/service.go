package chat

import "errors"

func QueryService(req *QueryReq) (*MessageResp, error) {
	if req == nil || req.MessageId <= 0 {
		return nil, errors.New("messageId required")
	}

	return &MessageResp{
		MessageId: req.MessageId,
		FromUid:   10001,
		Content:   "hello",
	}, nil
}

func CreateService(req *CreateReq, fromUid int64) (*MessageResp, error) {
	if req == nil || req.ToUid <= 0 || req.Content == "" {
		return nil, errors.New("toUid and content required")
	}

	return &MessageResp{
		MessageId: 1,
		FromUid:   fromUid,
		Content:   req.Content,
	}, nil
}

func ListService(req *ListReq) (*ListResp, error) {
	limit := 20
	if req != nil && req.Limit > 0 {
		limit = req.Limit
	}
	if limit > 100 {
		limit = 100
	}

	return &ListResp{
		Messages: []*MessageResp{
			{MessageId: 1, FromUid: 10001, Content: "hello"},
		},
	}, nil
}
