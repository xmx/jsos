package jsvm

import (
	"io"
	"sync"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jzip"
)

func New(mods ...ModuleRegister) (Engineer, error) {
	prog, err := onceCompileBabel()
	if err != nil {
		return nil, err
	}

	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	// babel need
	logFunc := func(goja.FunctionCall) goja.Value { return nil }
	_ = vm.Set("console", map[string]func(goja.FunctionCall) goja.Value{
		"log":   logFunc,
		"error": logFunc,
		"warn":  logFunc,
	})

	if _, err = vm.RunProgram(prog); err != nil {
		return nil, err
	}
	var transformFunc goja.Callable
	babel := vm.Get("Babel")
	if err = vm.ExportTo(babel.ToObject(vm).Get("transform"), &transformFunc); err != nil {
		return nil, err
	}

	transform := func(code string, opts map[string]any) (string, error) {
		if value, exx := transformFunc(babel, vm.ToValue(code), vm.ToValue(opts)); exx != nil {
			return "", exx
		} else {
			return value.ToObject(vm).Get("code").String(), nil
		}
	}
	eng := &jsEngine{
		vm:        vm,
		transform: transform,
	}
	rqu := &require{
		eng:     eng,
		modules: make(map[string]goja.Value, 16),
		sources: make(map[string]goja.Value, 16),
	}
	eng.require = rqu
	_ = vm.Set("require", rqu.require)

	if err = RegisterModules(eng, mods); err != nil {
		return nil, err
	}

	return eng, nil
}

type jsEngine struct {
	vm *goja.Runtime
	// transform Babel.transform()
	transform func(code string, opts map[string]any) (string, error)

	require *require

	mutex  sync.Mutex
	finals []func() error
}

func (jse *jsEngine) Runtime() *goja.Runtime {
	return jse.vm
}

func (jse *jsEngine) RunString(code string) (goja.Value, error) {
	commonJS, err := jse.transform(code, map[string]any{"plugins": []string{"transform-modules-commonjs"}})
	if err != nil {
		return nil, err
	}

	return jse.vm.RunString(commonJS)
}

func (jse *jsEngine) RunProgram(pgm *goja.Program) (goja.Value, error) {
	return jse.vm.RunProgram(pgm)
}

func (jse *jsEngine) RegisterModule(name string, module any, override bool) bool {
	return jse.require.register(name, module, override)
}

func (jse *jsEngine) AddFinalizer(finals ...func() error) {
	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range finals {
		if final != nil {
			jse.finals = append(jse.finals, final)
		}
	}
}

func (jse *jsEngine) Kill(cause any) {
	jse.vm.Interrupt(cause)

	jse.mutex.Lock()
	defer jse.mutex.Unlock()

	for _, final := range jse.finals {
		_ = final()
	}
	jse.finals = nil
	jse.require.kill()
}

func (jse *jsEngine) ClearInterrupt() {
	jse.vm.ClearInterrupt()
}

func (jse *jsEngine) RunJZip(filename string) (goja.Value, error) {
	jz, err := jzip.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	app := jz.Manifest.Application
	mainPath := app.Main
	if mainPath != "" {
		mainPath = "main"
	}
	mainFile, err := jz.Open(mainPath + ".js")
	if err != nil {
		return nil, err
	}
	defer mainFile.Close()

	data, err := io.ReadAll(mainFile)
	if err != nil {
		return nil, err
	}
	jse.require.source = jz

	return jse.RunString(string(data))
}
