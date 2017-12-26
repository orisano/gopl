package ex19

func F() (s string) {
	defer func() {
		if p := recover(); p != nil {
			s = p.(string)
		}
	}()
	panic("hello")
}
