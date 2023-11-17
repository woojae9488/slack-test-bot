package util

func Validate(err error) {
	if err != nil {
		panic(err)
	}
}
