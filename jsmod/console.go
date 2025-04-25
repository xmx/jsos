package jsmod

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jsvm"
)

func NewConsole() jsvm.ModuleRegister {
	return &stdConsole{}
}

type stdConsole struct {
	eng jsvm.Engineer
}

func (sc *stdConsole) RegisterModule(eng jsvm.Engineer) error {
	sc.eng = eng
	fields := map[string]any{
		"log":   sc.writeStdout,
		"debug": sc.writeStdout,
		"info":  sc.writeStdout,
		"error": sc.writeStderr,
		"warn":  sc.writeStderr,
	}

	return eng.Runtime().Set("console", fields)
}

func (sc *stdConsole) writeStdout(call goja.FunctionCall) goja.Value {
	msg, err := sc.format(call)
	if err != nil {
		return sc.eng.Runtime().NewGoError(err)
	}
	stdout := sc.eng.Device().Stdout()
	if _, err = stdout.Write(msg); err != nil {
		return sc.eng.Runtime().NewGoError(err)
	}

	return goja.Undefined()
}

func (sc *stdConsole) writeStderr(call goja.FunctionCall) goja.Value {
	msg, err := sc.format(call)
	if err != nil {
		return sc.eng.Runtime().NewGoError(err)
	}
	stderr := sc.eng.Device().Stderr()
	if _, err = stderr.Write(msg); err != nil {
		return sc.eng.Runtime().NewGoError(err)
	}

	return goja.Undefined()
}

func (sc *stdConsole) format(call goja.FunctionCall) ([]byte, error) {
	buf := new(bytes.Buffer)
	for _, arg := range call.Arguments {
		if err := sc.parse(buf, arg); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (sc *stdConsole) parse(buf *bytes.Buffer, val goja.Value) error {
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
		return sc.reflectParse(buf, v)
	}

	return nil
}

func (*stdConsole) reflectParse(buf *bytes.Buffer, v any) error {
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
