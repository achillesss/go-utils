package but

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// but is short for beautiful

func combineColor(color ColorName, isBackgroundColor bool) outPutSet {
	if color >= COLOR_BLACK && color <= COLOR_WHITE {
		if isBackgroundColor {
			return COLOR_BACKGROUND + outPutSet(color)
		}
		return COLOR_FOREGROUND + outPutSet(color)
	}
	return 0
}

func newPrintControlSequence(cc controlCode, params ...string) string {
	return fmt.Sprintf(PRINTER_FORMAT, strings.Join(params, ";"), cc)
}

func defaultControlSequence() string {
	return newPrintControlSequence(CONTROL_M, "0")
}

func addSet(src string, controlSequence string) string {
	return controlSequence + src
}

func delSet(src string, controlSequence string) string {
	for strings.Contains(src, controlSequence) {
		src = strings.Trim(src, controlSequence)
	}
	return src
}

func (x *printer) delSet(isSuffix bool, controlSequence string) {
	if isSuffix {
		x.suffix = delSet(x.suffix, controlSequence)
		return
	}
	x.prefix = delSet(x.prefix, controlSequence)
}

func (x *printer) addSet(isSuffix bool, controlSequence string) *printer {
	x.delSet(isSuffix, controlSequence)
	if isSuffix {
		x.suffix = addSet(x.suffix, controlSequence)
	} else {
		x.prefix = addSet(x.prefix, controlSequence)
	}
	return x
}

func (x *printer) control(isSuffix bool, cc controlCode, params ...string) *printer {
	cs := newPrintControlSequence(cc, params...)
	return x.addSet(isSuffix, cs)
}

func (x *printer) finalSufix() {
	x.control(true, CONTROL_M, "0")
}

func (x *printer) Color(color ColorName, isBackgroundColor bool) Buter {
	c := combineColor(color, isBackgroundColor)
	if c > 0 {
		x.control(false, CONTROL_M, fmt.Sprintf("%d", c))
	}
	return x
}

func (x *printer) Show(sets ...outPutSet) Buter {
	var setStr []string
	for _, set := range sets {
		setStr = append(setStr, fmt.Sprintf("%d", set))
	}
	x.control(false, CONTROL_M, setStr...)
	return x
}

func cutReturns(str *string) int {
	var count int
	for strings.HasSuffix(*str, "\n") {
		*str = strings.TrimSuffix(*str, "\n")
		count++
	}
	return count
}

func (x *printer) OneLinePrint(isLastUpdate bool) {
	length := len(fmt.Sprintf(x.format, x.args...))
	x.control(false, CONTROL_CLEAR_END)
	if !isLastUpdate {
		cutReturns(&x.format)
		x.control(true, CONTROL_D, fmt.Sprintf("%d", length))
		x.control(false, CONTROL_HIDE)
	} else {
		x.control(false, CONTROL_SHOW)
	}
	x.Print()
}

func (x *printer) Print() {
	x.finalSufix()
	f, args := x.formating()
	x.p(x.w, f, args...)
}

func (x *printer) formating() (formation string, args []interface{}) {
	returnsCount := cutReturns(&x.format)
	var returns string

	for i := 0; i < returnsCount; i++ {
		returns += "\n"
	}

	x.format = x.prefix + x.format + x.suffix + returns

	return x.format, x.args
}

func (x *printer) String() string {
	x.finalSufix()
	f, args := x.formating()
	return fmt.Sprintf(f, args...)
}

func NewButer(w io.Writer, format string, args ...interface{}) Buter {
	if w == nil {
		w = os.Stdout
	}
	return &printer{w: w, p: func(w io.Writer, format string, args ...interface{}) { fmt.Fprintf(w, format, args...) }, format: format, args: args}
}
