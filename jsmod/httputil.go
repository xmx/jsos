package jsmod

import (
	"net/http/httputil"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jsvm"
)

func NewHTTPUtil() jsvm.ModuleRegister {
	return new(stdHTTPUtil)
}

type stdHTTPUtil struct{}

func (sh *stdHTTPUtil) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"Proxy": sh.newProxy,
	}
	eng.RegisterModule("http/httputil", vals, true)

	return nil
}

func (*stdHTTPUtil) newProxy(_ goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	pxy := new(httputil.ReverseProxy)
	return vm.ToValue(pxy).(*goja.Object)
}
