package jsvm

import (
	"io"

	"github.com/dop251/goja"
)

type Device interface {
	Stdout() io.Writer
	Stderr() io.Writer
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

type Engineer interface {
	Runtime() *goja.Runtime
	Device() Device
	RunJZip(filepath string) (goja.Value, error)
	RunScript(name, code string) (goja.Value, error)
	RunProgram(pgm *goja.Program) (goja.Value, error)
	RegisterModule(name string, module any, override bool) bool
	AddFinalizer(finals ...func() error)
	Kill(cause any)
}

type ModuleRegister interface {
	RegisterModule(eng Engineer) error
}

func RegisterModules(eng Engineer, mods []ModuleRegister) error {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		if err := mod.RegisterModule(eng); err != nil {
			return err
		}
	}

	return nil
}
