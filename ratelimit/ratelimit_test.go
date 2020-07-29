package rate

import (
	"fmt"
	"testing"
	"time"
)

func TestRatelimit(t *testing.T) {
	var rules = NewRules()
	rules.AddRule("get", time.Second, 6, 0)
	rules.AddRule("get", time.Millisecond*500, 4, 0)
	rules.AddRule("get", time.Millisecond*100, 1, 0)

	var slice = make([]string, 1000)
	for i := range slice {
		var err = rules.Call("User", "get", 0)
		if err != nil {
			fmt.Printf("call %d error: %+v\n", i, err)
		}
		time.Sleep(time.Millisecond * 50)
	}
}
