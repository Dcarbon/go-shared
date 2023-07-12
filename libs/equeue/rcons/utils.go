package rcons

// PanicError :
func panicError(msg string, err error) {
	if nil != err {
		panic(msg + " error: " + err.Error())
	}
}
