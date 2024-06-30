package errors

import (
	"github.com/fatih/color"
)

type handler interface {
	Handle(err error)
}

func Check(h handler, err error) {
	if err != nil {
		color.Set(color.FgRed)
		if customHandler, ok := h.(handler); ok {
			customHandler.Handle(err)
		} else {
			panic(err)
		}
	}
}
