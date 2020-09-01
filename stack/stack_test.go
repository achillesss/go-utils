package stack

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultStack(t *testing.T) {
	var h = GetStackHash()
	fmt.Printf("stack hash: %v\n", h)
	var count, s = QueryStack(h)
	fmt.Printf("count: %d\tstack:\n%s\n", count, s)

	count, s = QueryStack("")
	fmt.Printf("count: %d\tstack:\n%s\n", count, s)
}

func TestStoreStack(t *testing.T) {
	var stack = ``
	var x = make([]string, 2)
	var start = time.Now()
	for range x {
		var s = getStackHash([]byte(stack))
		fmt.Printf("hash: %v\n", s)
	}
	fmt.Printf("cost: %v", time.Since(start))
}
