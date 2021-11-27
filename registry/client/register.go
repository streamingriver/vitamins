package registry

import (
	"fmt"
	"sync"
	"time"
)

func New(target string, servicename, myhost, myport string) *Pinger {
	return &Pinger{
		config: &config{
			servicename: servicename,
			myhost:      myhost,
			myport:      myport,
		},
		mu:      new(sync.RWMutex),
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

type Pinger struct {
	fetcher Fetcher
	config  *config
	mu      *sync.RWMutex
	host    string
	running bool

	auth *auth
}

func (p *Pinger) Start() {
	url := fmt.Sprintf("%s/%s/%s/%s", p.host, p.config.servicename, p.config.myhost, p.config.myport)
	p.mu.Lock()
	p.running = true
	p.mu.Unlock()
	for {
		p.mu.RLock()
		if p.running == false {
			p.mu.RLock()
			return
		}
		p.mu.RUnlock()
		p.fetcher.Fetch(url)
		time.Sleep(1 * time.Second)
	}
}

func (p *Pinger) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.running = false
}
