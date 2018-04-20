package but

import "io"

type Buter interface {
	Color(color ColorName, isBackgroundColor bool) Buter
	Show(...outPutSet) Buter
	Print()
	OneLinePrint(bool)
	String() string
}

type ColorName int
type outPutSet int
type controlCode string

const (
	// 0(黑)、1(红)、2(绿)、 3(黄)、4(蓝)、5(洋红)、6(青)、7(白)

	COLOR_BLACK ColorName = iota
	COLOR_RED
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
	COLOR_MAGENTA
	COLOR_CYAN
	COLOR_WHITE

	// 显示：0(默认)、1(粗体/高亮)、22(非粗体)、4(单条下划线)、24(无下划线)、5(闪烁)、25(无闪烁)、7(反显、翻转前景色和背景色)、27(无反显)

	SET_DEFAULT      outPutSet = 0
	SET_BOLD         outPutSet = 1
	SET_UNDERLINE    outPutSet = 4
	SET_BLINK        outPutSet = 5
	SET_REVERSAL     outPutSet = 7
	SET_UNBOLD       outPutSet = 22
	SET_NO_UNDERLINE outPutSet = 24
	SET_NO_BLINK     outPutSet = 25
	SET_UNREVERSAL   outPutSet = 27
	COLOR_FOREGROUND outPutSet = 30
	COLOR_BACKGROUND outPutSet = 40

	CONTROL_M         controlCode = "m"  // default
	CONTROL_A         controlCode = "A"  // up
	CONTROL_B         controlCode = "B"  // down
	CONTROL_C         controlCode = "C"  // right
	CONTROL_D         controlCode = "D"  // left
	CONTROL_CLEAR     controlCode = "2J" // clear monitor
	CONTROL_CLEAR_END controlCode = "K"  // clear to the end of row
	CONTROL_SAVE      controlCode = "s"  // save position
	CONTROL_RECOVER   controlCode = "u"  // recover position
	CONTROL_HIDE      controlCode = "?25l"
	CONTROL_SHOW      controlCode = "?25h"
	PRINTER_FORMAT                = "\033[%s%s"
)

type printer struct {
	p      func(w io.Writer, format string, args ...interface{})
	format string
	args   []interface{}
	prefix string
	suffix string
	w      io.Writer
}
