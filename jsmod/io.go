package jsmod

import (
	"io"

	"github.com/xmx/jsos/jsvm"
)

func NewIO() jsvm.ModuleLoader {
	return new(stdIO)
}

type stdIO struct{}

func (std *stdIO) LoadModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"copy":    io.Copy,
		"copyN":   io.CopyN,
		"discard": io.Discard,
	}
	eng.RegisterModule("io", vals, true)

	return nil
}
