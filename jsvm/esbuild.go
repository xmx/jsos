package jsvm

import (
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/parser"
	"github.com/evanw/esbuild/pkg/api"
)

// Transform transpiles the input source string and strip types from it.
// this is done using esbuild
func Transform(name, code string) (string, error) {
	opts := api.TransformOptions{
		Loader:         api.LoaderJS,
		Sourcefile:     name,
		Target:         api.ESNext,
		Format:         api.FormatCommonJS,
		Sourcemap:      api.SourceMapNone,
		SourcesContent: api.SourcesContentInclude,
		LegalComments:  api.LegalCommentsNone,
		Platform:       api.PlatformNeutral,
		LogLevel:       api.LogLevelSilent,
		Charset:        api.CharsetUTF8,
	}

	ret := api.Transform(code, opts)
	if err := esbuildCheckError(ret); err != nil {
		return "", err
	}

	return string(ret.Code), nil
}

func esbuildCheckError(result api.TransformResult) error {
	if len(result.Errors) == 0 {
		return nil
	}

	msg := result.Errors[0]
	err := &parser.Error{Message: msg.Text}

	if msg.Location != nil {
		err.Position = file.Position{
			Filename: msg.Location.File,
			Line:     msg.Location.Line,
			Column:   msg.Location.Column,
		}
	}

	return err
}
