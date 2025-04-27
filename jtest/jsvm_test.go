package jtest_test

import (
	"os"
	"testing"

	"github.com/xmx/jsos/jsmod"
	"github.com/xmx/jsos/jsvm"
)

func TestJZip(t *testing.T) {
	mods := []jsvm.ModuleRegister{
		jsmod.NewConsole(),
		jsmod.NewContext(),
		jsmod.NewExec(),
		jsmod.NewIO(),
		jsmod.NewOS(),
		jsmod.NewRuntime(),
		jsmod.NewTime(),
		jsmod.NewHTTP(),
	}

	eng, err := jsvm.New(mods...)
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Kill("结束")

	_, err = eng.RunJZip("demo/main.zip")
	t.Log(err)
}

func TestRun(t *testing.T) {
	name := "srv.js"
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	mods := []jsvm.ModuleRegister{
		jsmod.NewConsole(),
		jsmod.NewContext(),
		jsmod.NewExec(),
		jsmod.NewIO(),
		jsmod.NewOS(),
		jsmod.NewRuntime(),
		jsmod.NewTime(),
		jsmod.NewHTTP(),
	}

	eng, err := jsvm.New(mods...)
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Kill("结束")

	eng.Stdout().Attach(os.Stdout)
	eng.Stderr().Attach(os.Stderr)
	_, err = eng.RunScript(name, string(data))
	t.Log(err)
}
