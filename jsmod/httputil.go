package jsmod

import (
	"net/http/httputil"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jsvm"
)

func NewHTTPUtil() jsvm.ModuleLoader {
	return new(stdHTTPUtil)
}

type stdHTTPUtil struct{}

func (sh *stdHTTPUtil) LoadModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"Proxy": sh.newProxy,
	}
	eng.RegisterModule("http/httputil", vals, true)

	return nil
}

func (*stdHTTPUtil) newProxy(_ goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	pxy := new(httputil.ReverseProxy)
	hp := &httputilProxy{pxy: pxy}
	obj := vm.ToValue(hp).(*goja.Object)
	_ = obj.Set("setRewrite", hp.setRewrite)

	return obj
}

type httputilProxy struct {
	pxy *httputil.ReverseProxy
}

func (hp *httputilProxy) setRewrite(rewrite func(*httputil.ProxyRequest)) {
	hp.pxy.Rewrite = rewrite
}
