package apiserver

type Listener interface {
	Listen()
	OnMessage([]byte)
}
