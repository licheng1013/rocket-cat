package user

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResp struct {
	Uid   int64  `json:"uid"`
	Token string `json:"token"`
}

type LogoutReq struct {
	Token string `json:"token"`
}

type InfoReq struct {
	Uid int64 `json:"uid"`
}

type InfoResp struct {
	Uid      int64  `json:"uid"`
	Nickname string `json:"nickname"`
}
