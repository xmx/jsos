package jsmod

import (
	"runtime"

	"github.com/xmx/jsos/jsvm"
)

func NewRuntime() jsvm.ModuleRegister {
	return new(stdRuntime)
}

type stdRuntime struct{}

func (s *stdRuntime) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"memStats":     s.memStats,
		"goos":         runtime.GOOS,
		"goarch":       runtime.GOARCH,
		"gc":           runtime.GC,
		"numCPU":       runtime.NumCPU,
		"numGoroutine": runtime.NumGoroutine,
		"numCgoCall":   runtime.NumCgoCall,
		"version":      runtime.Version,
	}
	eng.RegisterModule("runtime", vals, true)

	return nil
}

func (s *stdRuntime) memStats() *runtime.MemStats {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	return stats
}
