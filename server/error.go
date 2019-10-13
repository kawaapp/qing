package server

import (
	"fmt"
)

func typeError(name string) error {
	return fmt.Errorf("value type should be %s", name)
}