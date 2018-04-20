package log

import (
	"flag"
	"fmt"
	"testing"

	but "github.com/achillesss/go-utils/but4print"
)

func turnOnInfo() {
	*infoOn = true
}
func turnOffInfo() {
	*infoOn = false
}
func turnOnWarn() {
	*warnOn = true
}
func turnOffWarn() {
	*warnOn = false
}
func turnOnError() {
	*errOn = true
}
func turnOffError() {
	*errOn = false
}
func turnOnTime() {
	*timeOn = true
}
func turnOffTime() {
	*timeOn = false
}

func TestLog(t *testing.T) {
	flag.Parse()
	turnOnTime()
}

func TestFuncName(t *testing.T) {
	func() {
		fmt.Printf("function name: %v\n", FuncName())
	}()
}

func TestFuncNameN(t *testing.T) {
	func() {
		fmt.Printf("function name: %v\n", FuncNameN(0))
	}()
	func() {
		fmt.Printf("function name: %v\n", FuncNameN(1))
	}()
	func() {
		fmt.Printf("function name: %v\n", FuncNameNP(2))
	}()
	func() {
		fmt.Printf("function name: %v\n", FuncNameNP(3))
	}()
}

func TestInfof(t *testing.T) {
	turnOnInfo()
	SetInfoColor(but.COLOR_WHITE, false)
	SetInfoColor(but.COLOR_BLUE, true)
	Infof("hello,world!")
	println()

	SetInfoColor(but.COLOR_WHITE, true)
	SetInfoColor(but.COLOR_BLUE, false)
	Infof("hello,world!")
	println()
	turnOffInfo()
}

func TestInfofln(t *testing.T) {
	turnOnInfo()
	SetInfoColor(but.COLOR_WHITE, false)
	SetInfoColor(but.COLOR_BLUE, true)
	Infofln("hello,world!")
	println()

	SetInfoColor(but.COLOR_WHITE, true)
	SetInfoColor(but.COLOR_BLUE, false)
	Infofln("hello,world!")
	println()
	turnOffInfo()
}

func TestWarnf(t *testing.T) {
	turnOnWarn()
	SetInfoColor(but.COLOR_WHITE, false)
	SetInfoColor(but.COLOR_YELLOW, true)
	Warningf("hello,world!")
	println()

	SetInfoColor(but.COLOR_WHITE, true)
	SetInfoColor(but.COLOR_YELLOW, false)
	Warningf("hello,world!")
	println()
	turnOffWarn()
}

func TestWarningfln(t *testing.T) {
	turnOnWarn()
	SetWarnColor(but.COLOR_WHITE, false)
	SetWarnColor(but.COLOR_YELLOW, true)
	Warningfln("hello,world!")
	println()

	SetWarnColor(but.COLOR_WHITE, true)
	SetWarnColor(but.COLOR_YELLOW, false)
	Warningfln("hello,world!")
	println()
	turnOffWarn()
}

func TestErrorf(t *testing.T) {
	turnOnError()
	SetErrorColor(but.COLOR_WHITE, false)
	SetErrorColor(but.COLOR_RED, true)
	Errorf("hello,world!")
	println()

	SetErrorColor(but.COLOR_WHITE, true)
	SetErrorColor(but.COLOR_RED, false)
	Errorf("hello,world!")
	println()
	turnOffError()

}
func TestErrorfln(t *testing.T) {
	turnOnError()
	SetErrorColor(but.COLOR_WHITE, false)
	SetErrorColor(but.COLOR_RED, true)
	Errorfln("hello,world!")
	println()

	SetErrorColor(but.COLOR_WHITE, true)
	SetErrorColor(but.COLOR_RED, false)
	Errorfln("hello,world!")
	println()
	turnOffError()
}

func TestLogInline(t *testing.T) {
	L(true, "%s Failed. Error: %v. Resp: %#v", FuncName(), "ERROR", "RESP")
	println()
	L(false, "%s Failed. Error: %v. Resp: %#v", FuncName(), "ERROR", "RESP")
	println()
}

func TestLogln(t *testing.T) {
	Lln(true, "%s Failed. Error: %v. Resp: %#v", FuncName(), "ERROR", "RESP")
	println()
	Lln(false, "%s Failed. Error: %v. Resp: %#v", FuncName(), "ERROR", "RESP")
	println()
}

func TestFormatError(t *testing.T) {
	turnOnError()
	err1 := fmt.Errorf("error1")
	Errorfln("err1: %v", err1)
	e1 := FmtErr(&err1)

	err2 := fmt.Errorf("error2")
	Infofln("err2: %v", err2)
	e2 := FmtErrP(&err2)

	Errorfln("ERR: %v\t%v", err1, e1)
	Errorfln("ERR: %v\t%v", err2, e2)
	turnOffError()
}
