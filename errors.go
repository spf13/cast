package cast

import (
	"errors"
	"fmt"
	"sync"
)

var errNegativeNotAllowed = errors.New("unable to cast negative value")

type lazyError struct {
	format string
	args   []interface{}
	once   sync.Once
	result string
}

func (e *lazyError) Error() string {
	e.once.Do(func() {
		e.result = fmt.Sprintf(e.format, e.args...)
	})
	return e.result
}

func newError(format string, args ...interface{}) error {
	return &lazyError{
		format: format,
		args:   args,
	}
}
