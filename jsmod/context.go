package jsmod

import (
	"context"

	"github.com/xmx/jsos/jsvm"
)

func NewContext() jsvm.ModuleLoader {
	return new(stdContext)
}

type stdContext struct{}

func (*stdContext) LoadModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"background":   context.Background,
		"todo":         context.TODO,
		"withCancel":   context.WithCancel,
		"withTimeout":  context.WithTimeout,
		"withValue":    context.WithValue,
		"withDeadline": context.WithDeadline,
	}
	eng.RegisterModule("context", vals, true)

	return nil
}
