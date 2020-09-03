package log

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	but "github.com/achillesss/go-utils/but4print"
)

func in(value interface{}, src ...interface{}) (ok bool) {
	for i := range src {
		if ok = value == src[i]; ok {
			break
		}
	}
	return
}

func newLogAgent() *logAgent {
	return &logAgent{symble: slash, end: inline, printType: logInformation}
}

func (a *logAgent) isNil() bool {
	return a == nil
}

func (a *logAgent) setSymble(on bool) *logAgent {
	a.symble = slash
	if !on {
		a.symble = point
	}
	return a
}

func (a *logAgent) setSkip(s int) *logAgent {
	a.skip = s
	return a
}

func (a *logAgent) setEnd(e string) *logAgent {
	if e == newline {
		a.end = e
	}
	return a
}

func (a *logAgent) setPrintType(p string) *logAgent {
	if in(p, logWarning, logError, logLog) {
		a.printType = p
	}
	return a
}

func (a *logAgent) funcName() string {
	pc, _, _, _ := runtime.Caller(a.skip + 1)
	name := runtime.FuncForPC(pc).Name()
	strs := strings.Split(name, a.symble)
	return strs[len(strs)-1]

}

func (a *logAgent) callerLine() string {
	if _, file, line, ok := runtime.Caller(a.skip + 1); ok {
		fileName := strings.Split(file, slash)
		return strings.Join([]string{fileName[len(fileName)-1], strconv.Itoa(line)}, "_")
	}

	return "unkown"
}

func (a *logAgent) String(format string, arg ...interface{}) string {
	if _, file, line, ok := runtime.Caller(a.skip + 1); ok {
		fileName := strings.Split(file, slash)
		timeTag := ""

		if *timeOn {
			timeTag = " " + time.Now().UTC().Format(time.StampMilli)
		}

		arg = append([]interface{}{a.printType, fileName[len(fileName)-1], line, timeTag}, arg...)
		format = "[%s_%v_%v%s] " + format + a.end
		printer := but.NewButer(nil, format, arg...)

		var (
			setBold   bool
			foreColor but.ColorName = -1
			backColor but.ColorName = -1
		)

		switch a.printType {
		case logWarning:
			foreColor = but.COLOR_WHITE
			backColor = but.COLOR_YELLOW
			if warnForeColor != nil {
				foreColor = *warnForeColor
			}
			if warnBackColor != nil {
				backColor = *warnBackColor
			}
			setBold = true

		case logError:
			foreColor = but.COLOR_WHITE
			backColor = but.COLOR_RED
			if errorForeColor != nil {
				foreColor = *errorForeColor
			}
			if errorBackColor != nil {
				backColor = *errorBackColor
			}
			setBold = true
		case logInformation:
			if infoForeColor != nil {
				foreColor = *infoForeColor
			}
			if infoBackColor != nil {
				backColor = *infoBackColor
			}
		}

		printer.Color(foreColor, false).Color(backColor, true)

		if setBold {
			printer.Show(but.SET_BOLD)
		}

		return printer.String()
	}
	return ""
}

func (a *logAgent) print(format string, arg ...interface{}) {
	str := a.String(format, arg...)
	if str != "" {
		print(str)
	}
}

func (a *logAgent) formatErr(err *error) error {
	if err != nil && *err != nil {
		*err = fmt.Errorf("%s fail. Desc: %v", a.funcName(), *err)
		return *err
	}
	return nil
}
