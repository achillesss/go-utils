package gerr

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	var errMap = map[int64]string{
		0: "Fail",
		1: "Success",
		2: "Internal",
		3: "NotFound",
		4: "Repeat",
	}

	SetDefaultErrorNotFoundCode(-1)
	RegisterErrors(errMap)

	var err = New(4).Wrap(3, "User 3 not found").Wrap(0)
	fmt.Printf("%v\n", err)
}
