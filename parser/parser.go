package parser

import (
	"encoding/base64"
	"reflect"
	"strings"
	"sync"
)

func New() *Parser {
	p := &Parser{
		mu:  &sync.RWMutex{},
		fns: make(map[string]interface{}),
	}
	return p
}

type Parser struct {
	mu  *sync.RWMutex
	fns map[string]interface{}
}

func (p *Parser) Register(name string, fn interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.fns[name] = fn
}

func (p *Parser) Call(s string) {
	cmd := pw(&s)

	p.mu.RLock()
	fn, ok := p.fns[cmd]
	p.mu.RUnlock()
	if !ok {
		return
	}

	funcArgs := reflect.ValueOf(fn).Type().NumIn()
	funcValue := reflect.ValueOf(fn)

	v := make([]reflect.Value, funcArgs)
	for i := 0; i < funcArgs-1; i++ {
		v[i] = reflect.ValueOf(pw(&s))
	}
	if funcArgs > 0 {
		ds, err1 := base64.StdEncoding.DecodeString(strings.Trim(s, "\n\r\t "))

		if err1 == nil {
			v[funcArgs-1] = reflect.ValueOf(string(ds))
		} else {
			v[funcArgs-1] = reflect.ValueOf(s)
		}
	}
	funcValue.Call(v)
}

func pw(s *string) string {
	l, rt, i := 100, "", 0

	for _, v := range *s {
		if v != ' ' {
			rt = rt + string(v)
		} else {
			break
		}
		i++
		if i >= l {
			break
		}
	}
	if len(*s) <= i {
		*s = ""
		ds, err := base64.StdEncoding.DecodeString(strings.Trim(rt, "\n\r\t "))

		if err == nil {
			return string(ds)
		}
		return rt
	}
	*s = (*s)[i+1:]

	ds, err := base64.StdEncoding.DecodeString(strings.Trim(rt, "\n\r\t "))

	if err == nil {
		return string(ds)
	}
	return rt
}
