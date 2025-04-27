package jsmod

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jsvm"
)

func NewConsole() jsvm.ModuleRegister {
	return &stdConsole{}
}

type stdConsole struct{}

func (sc *stdConsole) RegisterModule(eng jsvm.Engineer) error {
	cm := consoleModule{eng: eng}
	vals := map[string]any{
		"log":   cm.stdout,
		"debug": cm.stdout,
		"info":  cm.stdout,
		"error": cm.stderr,
		"warn":  cm.stderr,
	}

	return eng.Runtime().Set("console", vals)
}

type consoleModule struct {
	eng jsvm.Engineer
}

func (cm consoleModule) stdout(call goja.FunctionCall) goja.Value {
	vm, w := cm.eng.Runtime(), cm.eng.Stdout()
	return cm.writeTo(w, call, vm)
}

func (cm consoleModule) stderr(call goja.FunctionCall) goja.Value {
	vm, w := cm.eng.Runtime(), cm.eng.Stderr()
	return cm.writeTo(w, call, vm)
}

func (cm consoleModule) writeTo(w io.Writer, call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	msg, err := cm.format(call)
	if err != nil {
		return vm.NewGoError(err)
	}
	if _, err = w.Write(msg); err != nil {
		return vm.NewGoError(err)
	}

	return goja.Undefined()
}

func (cm consoleModule) format(call goja.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, arg := range call.Arguments {
		if err := cm.parse(buf, arg); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (cm consoleModule) parse(buf *bytes.Buffer, val goja.Value) error {
	switch {
	case goja.IsUndefined(val), goja.IsNull(val):
		buf.WriteString(val.String())
		return nil
	}

	export := val.Export()
	switch v := export.(type) {
	case fmt.Stringer:
		buf.WriteString(v.String())
	case string:
		buf.WriteString(v)
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
	case float64:
		buf.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
	case bool:
		buf.WriteString(strconv.FormatBool(v))
	case []byte:
		str := base64.StdEncoding.EncodeToString(v)
		buf.WriteString(str)
	case func(goja.FunctionCall) goja.Value:
		buf.WriteString("<Function>")
	case goja.ArrayBuffer:
		bs := v.Bytes()
		str := base64.StdEncoding.EncodeToString(bs)
		buf.WriteString(str)
	default:
		return cm.reflectParse(buf, v)
	}

	return nil
}

func (consoleModule) reflectParse(buf *bytes.Buffer, v any) error {
	vof := reflect.ValueOf(v)
	switch vof.Kind() {
	case reflect.String:
		buf.WriteString(vof.String())
	case reflect.Int64:
		buf.WriteString(strconv.FormatInt(vof.Int(), 10))
	case reflect.Float64:
		buf.WriteString(strconv.FormatFloat(vof.Float(), 'g', -1, 64))
	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(vof.Bool()))
	default:
		tmp := new(bytes.Buffer)
		if err := json.NewEncoder(tmp).Encode(v); err == nil && tmp.Len() != 0 {
			_, _ = buf.ReadFrom(tmp)
			return nil
		}
		vts := vof.Type().String()
		buf.WriteString(vts)
	}

	return nil
}
