package jsvm

import (
	"github.com/dop251/goja"
	"github.com/xmx/jsos/multio"
)

type Device interface {
	Stdout() multio.Writer
	Stderr() multio.Writer
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
