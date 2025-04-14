package jtest_test

import (
	"os"
	"testing"

	"github.com/xmx/jsos/jsmod"
	"github.com/xmx/jsos/jsvm"
)

func TestJZip(t *testing.T) {
	mods := []jsvm.ModuleRegister{
		jsmod.NewConsole(os.Stdout, os.Stderr),
		jsmod.NewContext(),
		jsmod.NewExec(),
		jsmod.NewIO(),
		jsmod.NewOS(),
		jsmod.NewRuntime(),
		jsmod.NewTime(),
	}
	eng, err := jsvm.New(mods...)
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Kill("结束")

	_, err = eng.RunJZip("demo/demo.zip")
	t.Log(err)
}
