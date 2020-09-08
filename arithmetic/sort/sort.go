package gosort

import (
	"fmt"
	"sort"
)

type Heaper interface {
	Pop() interface{}
	Push(interface{})
}

type HeapSorter interface {
	Heaper
	sort.Interface
}

type Debuger interface {
	Index(...int) string
}

type DebugSorter interface {
	Debuger
	sort.Interface
}

type DebugHeapSorter interface {
	Debuger
	Heaper
	sort.Interface
}

// sort float64
type float64Sorter []float64

// sort.Interface 接口
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
func SortFloat64(src []float64, f func(s sort.Interface)) {
	var s = float64Sorter(src)
	f(s)
}

// sort int64
type int64Sorter []int64

func (s int64Sorter) Len() int           { return len(s) }
func (s int64Sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s int64Sorter) Less(i, j int) bool { return s[i] < s[j] }

func SortInt64(src []int64, f func(sort.Interface)) {
	var s = int64Sorter(src)
	f(s)
}
