package events

import (
	"github.com/labstack/echo"
	"testing"
	"time"
)

func TestSubPub(t *testing.T) {
	bus := NewEventBus()

	results := make(chan string)

	bus.Subscribe("hello", func(c echo.Context, in interface{}) error {
		results <- "hello"
		return nil
	})
	bus.Subscribe("world", func(c echo.Context, in interface{}) error {
		results <- "world"
		return nil
	})

	go bus.Dispatch("hello", nil, nil)
	go bus.Dispatch("world", nil, nil)

	get := make(map[string]bool, 0)

Exit:
	for {
		select {
		case m := <-results:
			get[m] = true
			if get["hello"] && get["world"] {
				break Exit
			}
		case <-time.After(3 * time.Second):
			t.Error("timeout..")
			break Exit
		}
	}
}

func TestPanic(t *testing.T) {
	bus := NewEventBus()
	bus.Subscribe("hello", func(c echo.Context, in interface{}) error {
		panic("a test...")
	})
	bus.Dispatch("hello", nil, nil)
}
