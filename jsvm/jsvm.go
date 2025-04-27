package jsvm

import (
	"io"
	"slices"
	"sync"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jzip"
	"github.com/xmx/jsos/multio"
)

func New(mods ...ModuleLoader) (Engineer, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	eng := &jsEngine{
		vm:     vm,
		stdout: multio.New(),
		stderr: multio.New(),
	}
	if err := eng.enableRequire(); err != nil {
		return nil, err
	}
	if err := RegisterModules(eng, mods); err != nil {
		return nil, err
	}

	return eng, nil
}

type jsEngine struct {
	vm      *goja.Runtime
	require *require
	stdout  multio.Writer
	stderr  multio.Writer
	mutex   sync.Mutex
	finals  []func() error
}

func (jse *jsEngine) Runtime() *goja.Runtime {
	return jse.vm
}

func (jse *jsEngine) RunScript(name, code string) (goja.Value, error) {
	cjs, err := Transform(name, code)
	if err != nil {
		return nil, err
	}

	return jse.vm.RunScript(name, cjs)
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
	for _, final := range slices.Backward(jse.finals) {
		_ = final()
	}
	jse.finals = nil
}

func (jse *jsEngine) Stdout() multio.Writer {
	return jse.stdout
}

func (jse *jsEngine) Stderr() multio.Writer {
	return jse.stderr
}

func (jse *jsEngine) RunJZip(filename string) (goja.Value, error) {
	jz, err := jzip.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	app := jz.Manifest.Application
	mainName := app.Main
	if mainName == "" {
		mainName = "main"
	}

	mainPath := mainName + ".js"
	mainFile, err := jz.Open(mainPath)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer mainFile.Close()

	data, err := io.ReadAll(mainFile)
	if err != nil {
		return nil, err
	}
	jse.require.source = jz

	return jse.RunScript(mainPath, string(data))
}

func (jse *jsEngine) enableRequire() error {
	rqu := &require{
		eng:     jse,
		modules: make(map[string]goja.Value, 16),
		sources: make(map[string]goja.Value, 16),
	}
	jse.require = rqu
	if err := jse.vm.Set("require", rqu.require); err != nil {
		return err
	}
	jse.AddFinalizer(rqu.close)

	return nil
}
