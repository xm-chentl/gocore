package handles

import (
	"context"
)

type Handler interface {
	Call(context.Context) (interface{}, error)
}

type Handlers map[string]Handler

func (h Handlers) Has(key string) (ok bool) {
	_, ok = h[key]
	return
}

var (
	pool             = make(Handlers)
	methodOfHandlers = make(map[string][]Handler)
)

func Has(route string) bool {
	_, ok := pool[route]
	return ok
}

func Get(route string) Handler {
	return pool[route]
}

func Register(handlers Handlers) {
	if len(handlers) > 0 {
		for route, handler := range handlers {
			pool[route] = handler
		}
	}
}
