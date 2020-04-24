package funcaller

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/achillesss/go-utils/log"
)

func isFunction(f reflect.Value) bool {
	return f.Kind() == reflect.Func
}

func changeToFunction(function interface{}) (reflect.Value, bool) {
	fv := reflect.ValueOf(function)
	if isFunction(fv) {
		return fv, true
	}

	return fv, false
}

func changeToParams(params ...interface{}) []reflect.Value {
	var res []reflect.Value
	for _, p := range params {
		res = append(res, reflect.ValueOf(p))
	}
	return res
}

type FunctionCaller struct {
	f reflect.Value
	p []reflect.Value
}

func NewCaller(function interface{}, params ...interface{}) *FunctionCaller {
	var f FunctionCaller
	var ok bool
	f.f, ok = changeToFunction(function)
	if !ok {
		return nil
	}
	f.p = changeToParams(params...)
	return &f
}

func (c *FunctionCaller) Call(mustRecover bool, handlePanic func(err error)) []reflect.Value {
	if c == nil {
		return nil
	}

	var res []reflect.Value
	if mustRecover {
		defer func() {
			if r := recover(); r != nil {
				if handlePanic != nil {
					handlePanic(fmt.Errorf("%v", r))
				}
			}
		}()
	}

	res = c.f.Call(c.p)
	return res
}

func (c *FunctionCaller) GracefulRun(mustRecover bool, handlePanic func(error), beforeShutdown func()) {
	go c.Call(mustRecover, handlePanic)

	var signals = []os.Signal{
		os.Kill,
		os.Interrupt,    // terminal
		syscall.SIGTERM, // shutdown
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGHUP, // restart
		syscall.SIGUSR2,
	}

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	sig := <-signalChan
	switch sig {
	// 	case os.Kill:
	// 	case os.Interrupt: // terminal
	// 	case syscall.SIGTERM: // shutdown
	// 	case syscall.SIGINT:
	// 	case syscall.SIGQUIT:
	// 	case syscall.SIGHUP: // restart
	// 	case syscall.SIGUSR2:
	default:
		log.Infofln("ReceiveSig: %+v", sig)
		signal.Stop(signalChan)
		close(signalChan)
		beforeShutdown()
	}
}
