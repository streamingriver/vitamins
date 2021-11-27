package registry

import (
	"fmt"
	"sync"
	"time"
)

func New(schema string) *Registry {
	return &Registry{
		schema,
		make(map[string]*Item),
		&sync.RWMutex{},
		3,
	}
}

type Registry struct {
	schema   string
	registry map[string]*Item
	mu       *sync.RWMutex
	ttl      int64
}

type Item struct {
	Port string
	Host string
	Seen int64
}

func (r *Registry) SetTTL(ttl int64) {
	r.ttl = ttl
}

func (r Registry) GetURL(ch, file string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.registry[ch]

	if !ok {
		return "", &RegistryError{"Item not found", true, false}
	}

	if item.Seen <= time.Now().Unix() {
		return "", &RegistryError{"Item expired", false, true}
	}

	url := fmt.Sprintf("%s://%s:%s/%s%s", r.schema, item.Host, item.Port, ch, file)
	return url, nil
}

func (r *Registry) Ping(ch, host, port string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.registry[ch]
	if !ok {
		r.registry[ch] = &Item{
			Port: port,
			Host: host,
			Seen: time.Now().Unix() + r.ttl,
		}
	}
	r.registry[ch].Port = port
	r.registry[ch].Host = host
	r.registry[ch].Seen = time.Now().Unix() + r.ttl
}
