package user

import (
	"errors"
	"fmt"
)

func LoginService(req *LoginReq) (*LoginResp, error) {
	if req == nil || req.Account == "" || req.Password == "" {
		return nil, errors.New("account and password required")
	}

	return &LoginResp{
		Uid:   10001,
		Token: fmt.Sprintf("token-%s", req.Account),
	}, nil
}

func LogoutService(req *LogoutReq) error {
	if req == nil || req.Token == "" {
		return errors.New("token required")
	}
	return nil
}

func InfoService(req *InfoReq) (*InfoResp, error) {
	if req == nil || req.Uid <= 0 {
		return nil, errors.New("uid required")
	}

	return &InfoResp{
		Uid:      req.Uid,
		Nickname: "RocketCat",
	}, nil
}
