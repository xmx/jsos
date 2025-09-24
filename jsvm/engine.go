package jsvm

import (
	"context"

	"github.com/grafana/sobek"
	"github.com/xmx/jsos/multio"
)

type Engineer interface {
	Runtime() *sobek.Runtime
	RunScript(name, code string) (sobek.Value, error)
	RunProgram(pgm *sobek.Program) (sobek.Value, error)
	Finalizer() Finalizer
	Require() Requirer
	Output() (stdout, stderr multio.Writer)
	Kill(cause any)
	Context() context.Context
}

func New(parent context.Context) Engineer {
	vm := sobek.New()
	vm.SetFieldNameMapper(newJSONTagName())

	if parent == nil {
		parent = context.Background()
	}
	ctx, cancel := context.WithCancel(parent)

	svm := &sobekVM{
		vm:        vm,
		stdout:    multio.New(),
		stderr:    multio.New(),
		finalizer: newFinalizer(),
		ctx:       ctx,
		cancel:    cancel,
	}
	require := injectRequire(svm)
	svm.require = require
	context.AfterFunc(ctx, func() {
		svm.interrupt(context.Canceled)
	})

	return svm
}

type sobekVM struct {
	vm        *sobek.Runtime
	stdout    multio.Writer
	stderr    multio.Writer
	finalizer Finalizer
	require   *sobekRequire
	ctx       context.Context
	cancel    context.CancelFunc
}

func (svm *sobekVM) Runtime() *sobek.Runtime {
	return svm.vm
}

func (svm *sobekVM) RunScript(name, code string) (sobek.Value, error) {
	cjs, err := Transform(name, code)
	if err != nil {
		return nil, err
	}
	pgm, err := sobek.Compile(name, cjs, false)
	if err != nil {
		return nil, err
	}

	return svm.RunProgram(pgm)
}

func (svm *sobekVM) RunProgram(pgm *sobek.Program) (sobek.Value, error) {
	return svm.vm.RunProgram(pgm)
}

func (svm *sobekVM) Finalizer() Finalizer {
	return svm.finalizer
}

func (svm *sobekVM) Require() Requirer {
	return svm.require
}

func (svm *sobekVM) Output() (multio.Writer, multio.Writer) {
	return svm.stdout, svm.stderr
}

func (svm *sobekVM) Kill(cause any) {
	svm.interrupt(cause)
	svm.cancel()
}

func (svm *sobekVM) interrupt(cause any) {
	svm.finalizer.finalize()
	svm.vm.Interrupt(cause)
}

func (svm *sobekVM) Context() context.Context {
	return svm.ctx
}
