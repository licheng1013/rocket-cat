package common

// 断言工具

func AssertErr(err error) {
	if err != nil {
		print(err)
	}
}

const (
	GatewayName  = "gateway-game"
	ServicerName = "service-game"
)
