package apiserver

type Listener interface {
	Listen()
	OnMessage([]byte)
}

type Caller interface {
	Call(string)
}
