package core

// 断言工具

func AssertErr(err error) {
	if err != nil {
		print(err)
	}
}
