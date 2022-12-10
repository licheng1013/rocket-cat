package core

func AssertErr(err error) {
	if err != nil {
		print(err)
	}
}
