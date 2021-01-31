package api

func Must(err interface{}) {
	errVal, ok := err.(error)
	if ok {
		if errVal != nil {
			panic(err)
		}
		return
	}

	boolVal, ok := err.(bool)
	if ok {
		if boolVal {
			panic(err)
		}
	}
}
