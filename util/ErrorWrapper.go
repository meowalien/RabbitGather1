package util

import (
	"fmt"
)

// ErrorWrapper is a tool to deal with Nested error
type ErrorWrapper struct {
	Err error
}

func (w *ErrorWrapper) GetError() error {
	return w.Err
}

func (w *ErrorWrapper) Wrap(e error) {
	if e == nil {
		return
	}
	if w.Err != nil {
		w.Err = fmt.Errorf("%s -> %w", w.Err.Error(), e)
	} else {
		w.Err = e
	}
}
