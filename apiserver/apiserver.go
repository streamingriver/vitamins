package apiserver

import "gitlab.com/avarf/getenvs"

// New api server
func New(listeners ...Listener) *ApiServer {
	as := &ApiServer{}
	as.listeners = listeners
	return as
}

func NewDefault(callback Caller) *ApiServer {
	as := &ApiServer{}
	as.listeners = []Listener{
		&NatsListener{
			getenvs.GetEnvString("SERVICE_NAME", "no-name"),
			getenvs.GetEnvString("NATS_TOPIC", "configs"),
			getenvs.GetEnvString("NATS_URL", "nats://localhost:4222"),
			getenvs.GetEnvString("NATS_TOKEN", ""), callback,
		},
		&HttpListener{
			getenvs.GetEnvString("API_PORT", ":3080"), callback,
		},
	}
	return as
}

type ApiServer struct {
	listeners []Listener
}

func (as *ApiServer) Listen() {
	for _, l := range as.listeners {
		go l.Listen()
	}
}
