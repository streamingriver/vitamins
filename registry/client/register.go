package registry

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

func New(target string, servicename, myhost, myport, me string) *Pinger {
	return &Pinger{
		config: &config{
			servicename: servicename,
			myhost:      myhost,
			myport:      myport,
			me:          me,
		},
		delay:   1,
		host:    target,
		fetcher: &Fetch{},
		cb:      nil,
	}
}

type config struct {
	servicename string
	myhost      string
	myport      string
	me          string
}

type auth struct {
	username string
	password string
}

// Pinger struct
type Pinger struct {
	fetcher Fetcher
	config  *config
	host    string
	running int32

	auth *auth

	delay uint32

	cb func() string
}

// SetDelay between requests to registry
func (p *Pinger) SetDelay(d uint32) {
	atomic.StoreUint32(&p.delay, d)
}

func (p *Pinger) SetParamsFunc(f func() string) {
	p.cb = f
}

func (p *Pinger) SetAuth(username, password string) {
	p.fetcher.SetAuth(username, password)
}

// Start worker for periodic pings
func (p *Pinger) Start() {
	url := fmt.Sprintf("%s/%s/%s/%s", p.host, p.config.servicename, p.config.me, p.config.myport)
	atomic.StoreInt32(&p.running, 1)
	for {
		if atomic.LoadInt32(&p.running) == 0 {
			return
		}
		params := ""
		if p.cb != nil {
			params = "?" + p.cb()
		}
		err := p.fetcher.Fetch(url + params)
		if err != nil {
			log.Printf("%v", err)
		}
		time.Sleep(time.Duration(atomic.LoadUint32(&p.delay)) * time.Second)
	}
}

// Stop pinger worker
func (p *Pinger) Stop() {
	atomic.StoreInt32(&p.running, 0)
}
