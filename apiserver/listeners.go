package apiserver

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

type NatsMessage struct {
	Name string `json:"name"`
	// Data *json.RawMessage `json:"data"`
	Data string `json:"data"`
}

type NatsListener struct {
	Name     string
	Topic    string
	URL      string
	Token    string
	Callback func([]byte)
}

func (nl *NatsListener) Listen() {
	log.Printf("Connecting to nats " + nl.URL)
	nc, err := nats.Connect(nl.URL, nats.Token(nl.Token))
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	nc.Subscribe(nl.Topic, func(msg *nats.Msg) {
		nl.OnMessage(msg.Data)
	})
}

func (nl *NatsListener) OnMessage(b []byte) {
	if nl.Callback != nil {
		var natsMessage NatsMessage
		err := json.Unmarshal(b, &natsMessage)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		if natsMessage.Name != nl.Name {
			return
		}
		nl.Callback([]byte(natsMessage.Data))
	}
}

type HttpListener struct {
	Addr     string
	Callback func([]byte)
}

func (hl *HttpListener) Listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		hl.OnMessage(b)
	})
	log.Printf("Starting http listener on " + hl.Addr)
	log.Fatal(http.ListenAndServe(hl.Addr, nil))
}

func (hl *HttpListener) OnMessage(b []byte) {
	if hl.Callback != nil {
		hl.Callback(b)
	}
}
