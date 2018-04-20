package but

import (
	"testing"
	"time"
)

func TestButer(t *testing.T) {
	format0 := "%s\t\t"
	format1 := "Now:\t%s\n\n\n"
	arg0 := "Hello, world!"

	arg1 := func() string {
		return time.Now().String()
	}

	NewButer(
		nil,
		format0,
		arg0,
	).
		Color(COLOR_CYAN, false).
		Show(SET_BOLD).
		Print()

	NewButer(nil, format1, arg1()).
		Color(COLOR_RED, true).Color(COLOR_BLACK, false).
		Show(SET_UNDERLINE).
		Print()

	for i := 0; i < 101; i++ {
		NewButer(nil, "更新: %3d\n", i).Color(COLOR_BLUE, false).Show(SET_BOLD).OneLinePrint(i == 100)
		time.Sleep(time.Millisecond * 100)
	}
}
