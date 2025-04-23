package jsvm

import (
	"io"
	"sync"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jzip"
)

func New(mods ...ModuleRegister) (Engineer, error) {
	vm := goja.New()
	vm.SetFieldNameMapper(newFieldNameMapper("json"))
	eng := &jsEngine{
		vm:     vm,
		device: new(jsDevice),
	}
	rqu := &require{
		eng:     eng,
		modules: make(map[string]goja.Value, 16),
		sources: make(map[string]goja.Value, 16),
	}
	eng.require = rqu
	_ = vm.Set("require", rqu.require)

	if err := RegisterModules(eng, mods); err != nil {
		return nil, err
	}

	return eng, nil
}

type jsDevice struct {
	stdout io.Writer
	stderr io.Writer
}

func (jd *jsDevice) Stdout() io.Writer {
	if w := jd.stdout; w != nil {
		return w
	}

	return io.Discard
}

func (jd *jsDevice) Stderr() io.Writer {
	if w := jd.stderr; w != nil {
		return w
	}

	return io.Discard
}

func (jd *jsDevice) SetStdout(w io.Writer) {
	jd.stdout = w
}

func (jd *jsDevice) SetStderr(w io.Writer) {
	jd.stderr = w
}

type jsEngine struct {
	vm      *goja.Runtime
	device  *jsDevice
	require *require
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

func (jse *jsEngine) Device() Device {
	return jse.device
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
