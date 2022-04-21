package main

import (
	"fmt"
	"runtime/debug"
)

type any = interface{}

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]any
}

func wrapError(err error, messagef string, msgArgs ...any) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...),
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]any),
	}
}

func (err MyError) Error() string {
	return err.Message
}
