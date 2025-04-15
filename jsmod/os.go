package jsmod

import (
	"os"

	"github.com/xmx/jsos/jsvm"
)

func NewOS() jsvm.ModuleRegister {
	return new(stdOS)
}

type stdOS struct{}

func (std *stdOS) RegisterModule(eng jsvm.Engineer) error {
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
	}
	eng.RegisterModule("os", vals, true)

	return nil
}
