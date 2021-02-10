package lease

import (
	"errors"
	"fmt"
)

type internalDiscardLogger struct {}

func (internalDiscardLogger) Alert(a ...interface{}) error {
	return errors.New(fmt.Sprint(a...))
}

func (internalDiscardLogger) Debug(...interface{}) {
	// Nothing here.
}

func (internalDiscardLogger) Error(a ...interface{}) error {
	return errors.New(fmt.Sprint(a...))
}

func (internalDiscardLogger) Warn(...interface{}) {
	// Nothing here.
}
