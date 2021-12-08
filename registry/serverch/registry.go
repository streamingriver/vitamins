package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

func New(schema string) *Registry {
	return &Registry{
		schema,
		make(map[string]int),
		make(map[string][]*Item),
		&sync.RWMutex{},
		5,
		false,
	}
}

type Registry struct {
	schema     string
	nodes      map[string]int
	registry   map[string][]*Item
	mu         *sync.RWMutex
	ttl        int64
	dockerMode bool
}

type Item struct {
	Port string
	Host string
	Name string
	Seen int64
}

func (r *Registry) SetTTL(ttl int64) {
	r.ttl = ttl
}

func (r *Registry) SetDockerMode(mode bool) {
	r.dockerMode = mode
}

func (r Registry) GetURL(key string, name, file string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	items, ok := r.registry[name]

	if !ok {
		return "", &RegistryError{"Item not found", true, false}
	}

	nodesCount := r.nodes[name]
	var item *Item

	djbIndex := int(djb33(djb33Seed, key)) % r.nodes[name]

	item = items[0]

	if nodesCount > 1 {
		item = items[djbIndex]
	}
	tries := 0

	for {
		if item.Seen >= time.Now().Unix() {
			break
		}
		if tries >= nodesCount {
			break
		}
		item = items[(djbIndex+tries+1)%r.nodes[name]]
		tries++
	}

	if item.Seen <= time.Now().Unix() {
		return "", &RegistryError{"Item expired or no alive nodes found", false, true}
	}

	url := fmt.Sprintf("%s://%s:%s/%s/%s", r.schema, item.Host, item.Port, name, file)
	return url, nil
}

func (r *Registry) Ping(ch, host, port string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.registry[ch]
	if !ok {
		r.nodes[ch] = 0
	}
	itemIndex := r.hostPortExists(ch, host, port)
	if !ok || itemIndex == -1 {
		r.registry[ch] = append(r.registry[ch], &Item{
			Port: port,
			Host: host,
			Seen: time.Now().Unix() + r.ttl,
			Name: ch,
		})
		r.nodes[ch] = r.nodes[ch] + 1

	}
	if r.dockerMode {
		itemIndexForUpdate := r.hostPortExists(ch, host, "")
		if itemIndexForUpdate != -1 {
			r.registry[ch][itemIndexForUpdate].Port = port
		}
	}

	itemIndex = r.hostPortExists(ch, host, port)
	r.registry[ch][itemIndex].Seen = time.Now().Unix() + r.ttl
}

func (r *Registry) Debug() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, _ := json.MarshalIndent(r.registry, "", " ")
	log.Printf("%s", b)
}

func (r *Registry) hostPortExists(name, host, port string) int {

	for idx, item := range r.registry[name] {
		if item.Host == host {
			if item.Port == port && port != "" {
				return idx
			} else if port == "" {
				return idx
			}
		}
	}
	return -1
}
