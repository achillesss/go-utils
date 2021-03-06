package log

// Author: Achillesss
// date: 2017-03-10 23:12:00

import (
	"flag"

	but "github.com/achillesss/go-utils/but4print"
)

var (
	infoOn = flag.Bool("info", false, "whether print 'info', default off")
	warnOn = flag.Bool("warn", false, "whether print 'warning', default off")
	errOn  = flag.Bool("err", true, "whether print 'error', default on")
	timeOn = flag.Bool("time", false, "whether print with a time.Now().UTC().Format(time.RFC3339)[:19] tag")

	infoForeColor *but.ColorName
	infoBackColor *but.ColorName

	warnForeColor *but.ColorName
	warnBackColor *but.ColorName

	errorForeColor *but.ColorName
	errorBackColor *but.ColorName
)

// Parse parses flags

func SetInfoColor(c but.ColorName, isBackgroundColor bool) {
	if isBackgroundColor {
		infoBackColor = &c
		return
	}
	infoForeColor = &c
}

func SetWarnColor(c but.ColorName, isBackgroundColor bool) {
	if isBackgroundColor {
		warnBackColor = &c
		return
	}
	warnForeColor = &c
}

func SetErrorColor(c but.ColorName, isBackgroundColor bool) {
	if isBackgroundColor {
		errorBackColor = &c
		return
	}
	errorForeColor = &c
}

type formatErrCover func(skip int, pkgName bool, err *error) error

func formatErr(skip int, pkgName bool, err *error) error {
	return newLogAgent().setSkip(skip + 1).setSymble(pkgName).formatErr(err)
}

func print_(ok bool, skip int, printType, end string, format string, arg ...interface{}) {
	if ok {
		newLogAgent().setSkip(skip+2).setPrintType(printType).setEnd(end).print(format, arg...)
	}
}

func format(skip int, printType, end string, fmt string, arg ...interface{}) string {
	return newLogAgent().setSkip(skip+1).setPrintType(printType).setEnd(end).String(fmt, arg...)
}

func funcName(skip int, on bool) string {
	return newLogAgent().setSkip(skip + 1).setSymble(on).funcName()
}

func callerLine(skip int, on bool) string {
	return newLogAgent().setSkip(skip + 1).setSymble(on).callerLine()
}

// FuncName returns name of the function which calls it
func FuncName() string {
	return funcName(1, false)
}

// FuncNameP returns name of the function which calls it with package name
func FuncNameP() string {
	return funcName(1, true)
}

// FuncNameN returns name of function skipped by n
func FuncNameN(skip int) string {
	return funcName(skip+1, false)
}

// FuncNameNP returns name of function with package name by n
func FuncNameNP(skip int) string {
	return funcName(skip+1, true)
}

func formatErrEntrance(skip int, pkgName bool, err *error) error {
	if err != nil && *err != nil {
		return formatErr(skip+2, pkgName, err)
	}
	return nil
}

// FmtErr formats an error with name of the function which calls it
func FmtErr(err *error) error {
	return formatErrEntrance(1, false, err)
}

// FmtErrP differs from FmtErr in that it formats an error WITH PACKAGE NAME IN ADDITION
func FmtErrP(err *error) error {
	return formatErrEntrance(1, true, err)
}

// FmtErrN formats an error with name of the function which calls it skipped by skip
func FmtErrN(skip int, err *error) error {
	return formatErrEntrance(skip+1, false, err)
}

// FmtErrNP differs from FmtErrN in that it formats an error WITH PACKAGE NAME IN ADDITION
func FmtErrNP(skip int, err *error) error {
	return formatErrEntrance(skip+1, true, err)
}

// L logs a description when a function response is not ok
func L(ok bool, format string, arg ...interface{}) {
	print_(ok, 1, logLog, "", format, arg...)
}

// Lln differs from L in that it create a newline after loging
func Lln(ok bool, format string, arg ...interface{}) {
	print_(ok, 1, logLog, newline, format, arg...)
}

// Lln differs from L in that it create a newline after loging with a skip
func LlnN(ok bool, skip int, format string, arg ...interface{}) {
	print_(ok, 1+skip, logLog, newline, format, arg...)
}

// Infof prints information inline
func Infof(format string, arg ...interface{}) {
	print_(*infoOn, 1, "", "", format, arg...)
}

// Infofln prints information and create new line
func Infofln(format string, arg ...interface{}) {
	print_(*infoOn, 1, "", newline, format, arg...)
}

// InfoflnN prints information and create new line with a skip
func InfoflnN(skip int, format string, arg ...interface{}) {
	print_(*infoOn, 1+skip, "", newline, format, arg...)
}

// Warningf prints information inline
func Warningf(format string, arg ...interface{}) {
	print_(*warnOn || *infoOn, 1, logWarning, "", format, arg...)
}

// Warningfln prints information and create new line
func Warningfln(format string, arg ...interface{}) {
	print_(*warnOn || *infoOn, 1, logWarning, newline, format, arg...)
}

// WarningflnN prints information and create new line with a skip
func WarningflnN(skip int, format string, arg ...interface{}) {
	print_(*warnOn || *infoOn, 1+skip, logWarning, newline, format, arg...)
}

// Errorf prints information inline
func Errorf(format string, arg ...interface{}) {
	print_(*errOn || *warnOn || *infoOn, 1, logError, "", format, arg...)
}

// Errorfln prints information and create new line
func Errorfln(format string, arg ...interface{}) {
	print_(*errOn || *warnOn || *infoOn, 1, logError, newline, format, arg...)
}

// ErrorflnN prints information and create new line with a skip
func ErrorflnN(skip int, format string, arg ...interface{}) {
	print_(*errOn || *warnOn || *infoOn, 1+skip, logError, newline, format, arg...)
}

func FormatInfofln(fmt string, args ...interface{}) string {
	return format(1, "", newline, fmt, args...)
}

func FormatWarningfln(fmt string, args ...interface{}) string {
	return format(1, logWarning, newline, fmt, args...)
}

func FormatErrorfln(fmt string, args ...interface{}) string {
	return format(1, logError, newline, fmt, args...)
}

func CallerLine(skip int) string {
	return callerLine(skip+1, false)
}
