package funcaller

import (
	"fmt"
	"reflect"

	"runtime/debug"

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

func (c *FunctionCaller) Call(mustRecover bool) []reflect.Value {
	if c == nil {
		return nil
	}

	if mustRecover {
		defer func() {
			if r := recover(); r != nil {
				log.Errorfln("Panic", fmt.Errorf("panic: %v\r\n\r\n%s", r, debug.Stack()))
			}
		}()
	}

	return c.f.Call(c.p)
}
