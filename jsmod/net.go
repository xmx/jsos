package jsmod

import (
	"net"

	"github.com/xmx/jsos/jsvm"
)

func NewNet() jsvm.ModuleRegister {
	return new(stdNet)
}

type stdNet struct{}

func (sn *stdNet) RegisterModule(eng jsvm.Engineer) error {
	vals := map[string]any{
		"splitHostPort": net.SplitHostPort,
		"joinHostPort":  net.JoinHostPort,
		"parseCIDR":     net.ParseCIDR,
		"parseIP":       net.ParseIP,
		"parseMAC":      net.ParseMAC,
	}
	eng.RegisterModule("net", vals, true)

	return nil
}
