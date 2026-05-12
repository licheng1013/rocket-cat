package user

// LoginReq 表示用户登录请求。
type LoginReq struct {
	Account  string `json:"account"`  // 登录账号
	Password string `json:"password"` // 登录密码
}

// LoginResp 表示用户登录响应。
type LoginResp struct {
	Uid   int64  `json:"uid"`   // 用户 ID
	Token string `json:"token"` // 登录令牌
}

// LogoutReq 表示用户退出请求。
type LogoutReq struct {
	Token string `json:"token"` // 登录令牌
}

// InfoReq 表示用户信息查询请求。
type InfoReq struct {
	Uid int64 `json:"uid"` // 用户 ID
}

// InfoResp 表示用户信息查询响应。
type InfoResp struct {
	Uid      int64  `json:"uid"`      // 用户 ID
	Nickname string `json:"nickname"` // 用户昵称
}
