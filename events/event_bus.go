package events

import (
	"github.com/labstack/echo"
	"log"
	"sync"
)

type Handler func(c echo.Context, in interface{}) error

type EventBus interface {
	Subscribe(name string, handler Handler)
	Dispatch(name string, c echo.Context, in interface{})
}

type eventBus struct {
	sync.Mutex
	consumers map[string][]Handler
}

func (bus *eventBus) Subscribe(name string, handler Handler) {
	bus.Lock()
	bus.consumers[name] = append(bus.consumers[name], handler)
	bus.Unlock()
}

func (bus *eventBus) Dispatch(name string, c echo.Context, in interface{}) {
	bus.Lock()
	handlers := bus.consumers[name]
	bus.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered in EventBus.Dispatch()", r)
		}
	}()

	// Make it Async TODO
	for _, handler := range handlers {
		if err := handler(c, in); err != nil {
			log.Println(err)
		}
	}
}

func Dispatch(name string, c echo.Context, in interface{}) {
	defaultBus.Dispatch(name, c, in)
}

func Subscribe(name string, h Handler) {
	defaultBus.Subscribe(name, h)
}

func NewEventBus() EventBus {
	return &eventBus{
		consumers: make(map[string][]Handler),
	}
}

// default bus
var defaultBus EventBus

func init() {
	defaultBus = NewEventBus()
}
