package gosort

import (
	"fmt"
)

// Sorter
// 实现排序需要的三个方法
// 在其他排序方法中被调用
type Sorter interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

type Heaper interface {
	Pop() interface{}
	Push(interface{})
}

type HeapSorter interface {
	Heaper
	Sorter
}

type Debuger interface {
	Index(...int) string
}

type DebugSorter interface {
	Debuger
	Sorter
}

type DebugHeapSorter interface {
	Debuger
	Heaper
	Sorter
}

// sort float64
type float64Sorter []float64

// Sorter 接口
func (s float64Sorter) Len() int           { return len(s) }
func (s float64Sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s float64Sorter) Less(i, j int) bool { return s[i] < s[j] }

// Debuger 接口
func (s float64Sorter) Index(index ...int) string {
	var temp = make(float64Sorter, len(index))
	for i, j := range index {
		temp[i] = s[j]
	}
	return fmt.Sprintf("%+#v", temp)
}

// Heaper 接口
func (s *float64Sorter) Pop() interface{} {
	var old = *s
	var l = old.Len()
	var p = old[l-1]
	*s = old[:l-1]
	return p
}

func (s *float64Sorter) Push(p interface{}) {
	*s = append(*s, p.(float64))
}

// 外部 Sort 方法
func SortFloat64(src []float64, f func(Sorter)) {
	var s = float64Sorter(src)
	f(s)
}

// sort int64
type int64Sorter []int64

func (s int64Sorter) Len() int           { return len(s) }
func (s int64Sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64Sorter) Less(i, j int) bool { return s[i] < s[j] }

func SortInt64(src []int64, f func(Sorter)) {
	var s = int64Sorter(src)
	f(s)
}
