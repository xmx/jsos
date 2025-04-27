package jsmod

import (
	"errors"
	"net/http"

	"github.com/dop251/goja"
	"github.com/xmx/jsos/jsvm"
)

func NewHTTP() jsvm.ModuleLoader {
	return new(stdHTTP)
}

type stdHTTP struct{}

func (sh *stdHTTP) LoadModule(eng jsvm.Engineer) error {
	srv := &httpServer{eng: eng}
	vals := map[string]any{
		"statusContinue":           http.StatusContinue,
		"statusSwitchingProtocols": http.StatusSwitchingProtocols,
		"statusProcessing":         http.StatusProcessing,
		"statusEarlyHints":         http.StatusEarlyHints,

		"statusOK":                   http.StatusOK,
		"statusCreated":              http.StatusCreated,
		"statusAccepted":             http.StatusAccepted,
		"statusNonAuthoritativeInfo": http.StatusNonAuthoritativeInfo,
		"statusNoContent":            http.StatusNoContent,
		"statusResetContent":         http.StatusResetContent,
		"statusPartialContent":       http.StatusPartialContent,
		"statusMultiStatus":          http.StatusMultiStatus,
		"statusAlreadyReported":      http.StatusAlreadyReported,
		"statusIMUsed":               http.StatusIMUsed,

		"statusMultipleChoices":   http.StatusMultipleChoices,
		"statusMovedPermanently":  http.StatusMovedPermanently,
		"statusFound":             http.StatusFound,
		"statusSeeOther":          http.StatusSeeOther,
		"statusNotModified":       http.StatusNotModified,
		"statusUseProxy":          http.StatusUseProxy,
		"statusTemporaryRedirect": http.StatusTemporaryRedirect,
		"statusPermanentRedirect": http.StatusPermanentRedirect,

		"statusBadRequest":                   http.StatusBadRequest,
		"statusUnauthorized":                 http.StatusUnauthorized,
		"statusPaymentRequired":              http.StatusPaymentRequired,
		"statusForbidden":                    http.StatusForbidden,
		"statusNotFound":                     http.StatusNotFound,
		"statusMethodNotAllowed":             http.StatusMethodNotAllowed,
		"statusNotAcceptable":                http.StatusNotAcceptable,
		"statusProxyAuthRequired":            http.StatusProxyAuthRequired,
		"statusRequestTimeout":               http.StatusRequestTimeout,
		"statusConflict":                     http.StatusConflict,
		"statusGone":                         http.StatusGone,
		"statusLengthRequired":               http.StatusLengthRequired,
		"statusPreconditionFailed":           http.StatusPreconditionFailed,
		"statusRequestEntityTooLarge":        http.StatusRequestEntityTooLarge,
		"statusRequestURITooLong":            http.StatusRequestURITooLong,
		"statusUnsupportedMediaType":         http.StatusUnsupportedMediaType,
		"statusRequestedRangeNotSatisfiable": http.StatusRequestedRangeNotSatisfiable,
		"statusExpectationFailed":            http.StatusExpectationFailed,
		"statusTeapot":                       http.StatusTeapot,
		"statusMisdirectedRequest":           http.StatusMisdirectedRequest,
		"statusUnprocessableEntity":          http.StatusUnprocessableEntity,
		"statusLocked":                       http.StatusLocked,
		"statusFailedDependency":             http.StatusFailedDependency,
		"statusTooEarly":                     http.StatusTooEarly,
		"statusUpgradeRequired":              http.StatusUpgradeRequired,
		"statusPreconditionRequired":         http.StatusPreconditionRequired,
		"statusTooManyRequests":              http.StatusTooManyRequests,
		"statusRequestHeaderFieldsTooLarge":  http.StatusRequestHeaderFieldsTooLarge,
		"statusUnavailableForLegalReasons":   http.StatusUnavailableForLegalReasons,

		"statusInternalServerError":           http.StatusInternalServerError,
		"statusNotImplemented":                http.StatusNotImplemented,
		"statusBadGateway":                    http.StatusBadGateway,
		"statusServiceUnavailable":            http.StatusServiceUnavailable,
		"statusGatewayTimeout":                http.StatusGatewayTimeout,
		"statusHTTPVersionNotSupported":       http.StatusHTTPVersionNotSupported,
		"statusVariantAlsoNegotiates":         http.StatusVariantAlsoNegotiates,
		"statusInsufficientStorage":           http.StatusInsufficientStorage,
		"statusLoopDetected":                  http.StatusLoopDetected,
		"statusNotExtended":                   http.StatusNotExtended,
		"statusNetworkAuthenticationRequired": http.StatusNetworkAuthenticationRequired,

		"ServeMux":           sh.newServeMux,
		"listenAndServe":     srv.listenAndServe,
		"canonicalHeaderKey": http.CanonicalHeaderKey,
		"Client":             sh.newClient,
	}
	eng.RegisterModule("http", vals, true)

	return nil
}

type httpServer struct {
	eng jsvm.Engineer
}

func (hs httpServer) listenAndServe(addr string, handler http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: handler}
	hs.eng.AddFinalizer(srv.Close)
	err := srv.ListenAndServe()
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (*stdHTTP) newClient(_ goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	cli := &http.Client{}
	return vm.ToValue(cli).(*goja.Object)
}

func (*stdHTTP) newServeMux(_ goja.ConstructorCall, vm *goja.Runtime) *goja.Object {
	mux := http.NewServeMux()
	return vm.ToValue(mux).(*goja.Object)
}
