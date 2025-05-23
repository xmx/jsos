package jsvm

import (
	"github.com/dop251/goja"
	"github.com/xmx/jsos/multio"
)

type Engineer interface {
	Runtime() *goja.Runtime
	RunJZip(filepath string) (goja.Value, error)
	RunScript(name, code string) (goja.Value, error)
	RunProgram(pgm *goja.Program) (goja.Value, error)
	RegisterModule(name string, module any, override bool) bool
	AddFinalizer(finals ...func() error)
	Kill(cause any)

	Stdout() multio.Writer
	Stderr() multio.Writer
}

type ModuleLoader interface {
	LoadModule(eng Engineer) error
}

func RegisterModules(eng Engineer, mods []ModuleLoader) error {
	for _, mod := range mods {
		if mod == nil {
			continue
		}
		if err := mod.LoadModule(eng); err != nil {
			return err
		}
	}

	return nil
}
