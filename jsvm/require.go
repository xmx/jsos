package jsvm

import (
	"fmt"
	"io"
	"os"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jzip"
)

type require struct {
	eng     *jsEngine
	modules map[string]goja.Value
	sources map[string]goja.Value
	source  *jzip.JZip
}

func (rqu *require) register(name string, module any, override bool) bool {
	_, exists := rqu.modules[name]
	if exists && !override {
		return false
	}
	rqu.modules[name] = rqu.eng.vm.ToValue(module)

	return true
}

func (rqu *require) kill() {
	if rqu.source != nil {
		_ = rqu.source.Close()
	}
}

func (rqu *require) require(call goja.FunctionCall) goja.Value {
	name := call.Argument(0).String()
	val, exists := rqu.loadBootstrap(name)
	if exists {
		return val
	}

	var err error
	if val, exists, err = rqu.loadApplication(name); err == nil && exists {
		return val
	}

	vm := rqu.eng.Runtime()
	if !os.IsNotExist(err) {
		panic(vm.NewGoError(err))
	}

	msg := fmt.Sprintf("cannot find module '%s'", name)
	panic(vm.NewTypeError(msg))
}

func (rqu *require) loadBootstrap(name string) (goja.Value, bool) {
	val, exists := rqu.modules[name]
	return val, exists
}

func (rqu *require) loadApplication(name string) (goja.Value, bool, error) {
	if rqu.source == nil {
		return nil, false, nil
	}
	if val, exists := rqu.sources[name]; exists {
		return val, true, nil
	}

	filename := name + ".js"
	file, err := rqu.source.Open(filename)
	if err != nil {
		return nil, false, err
	}
	defer file.Close()
	code, err := io.ReadAll(file)
	if err != nil {
		return nil, false, err
	}

	vm := rqu.eng.Runtime()
	module := vm.NewObject()
	_ = vm.Set("module", module)
	if _, err = rqu.eng.RunScript(filename, string(code)); err != nil {
		return nil, false, err
	}
	exports := module.Get("exports").ToObject(vm)
	rqu.sources[name] = exports

	return exports, true, nil
}
