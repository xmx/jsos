package jsmod

import (
	"os"

	"github.com/xmx/jsos/jsvm"
)

func NewOS() jsvm.ModuleLoader {
	return new(stdOS)
}

type stdOS struct{}

func (std *stdOS) LoadModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"getpid":       os.Getpid,
		"open":         os.Open,
		"hostname":     os.Hostname,
		"tempDir":      os.TempDir,
		"getenv":       os.Getenv,
		"setenv":       os.Setenv,
		"unsetenv":     os.Unsetenv,
		"userCacheDir": os.UserCacheDir,
		"environ":      os.Environ,
		"expand":       os.Expand,
		"expandEnv":    os.ExpandEnv,
		"getuid":       os.Getuid,
		"geteuid":      os.Geteuid,
		"getgid":       os.Getgid,
		"getegid":      os.Getegid,
		"getgroups":    os.Getgroups,
		"getpagesize":  os.Getpagesize,
		"getppid":      os.Getppid,
		"getwd":        os.Getwd,
		"stdout":       os.Stdout,
		"stderr":       os.Stderr,
		"create":       os.Create,
	}
	eng.RegisterModule("os", vals, true)

	return nil
}
