package registry

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

func New(target string, servicename, myhost, myport string) *Pinger {
	return &Pinger{
		config: &config{
			servicename: servicename,
			myhost:      myhost,
			myport:      myport,
		},
		delay:   1,
		host:    target,
		fetcher: &Fetch{},
	}
}

type config struct {
	servicename string
	myhost      string
	myport      string
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
}

// SetDelay between requests to registry
func (p *Pinger) SetDelay(d uint32) {
	atomic.StoreUint32(&p.delay, d)
}

// Start worker for periodic pings
func (p *Pinger) Start() {
	url := fmt.Sprintf("%s/%s/%s/%s", p.host, p.config.servicename, p.config.myhost, p.config.myport)
	atomic.StoreInt32(&p.running, 1)
	for {
		if atomic.LoadInt32(&p.running) == 0 {
			return
		}
		err := p.fetcher.Fetch(url)
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
