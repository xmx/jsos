package jsmod

import (
	"net/url"

	"github.com/xmx/jsos/jsvm"
)

func NewURL() jsvm.ModuleRegister {
	return new(stdURL)
}

type stdURL struct{}

func (sn *stdURL) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"parse":           url.Parse,
		"joinPath":        url.JoinPath,
		"parseQuery":      url.ParseQuery,
		"parseRequestURI": url.ParseRequestURI,
		"pathEscape":      url.PathEscape,
		"pathUnescape":    url.PathUnescape,
		"queryEscape":     url.QueryEscape,
		"queryUnescape":   url.QueryUnescape,
		"userPassword":    url.UserPassword,
	}
	eng.RegisterModule("url", vals, true)

	return nil
}
