package gosort

import (
	"fmt"
	"reflect"
)

// Sorter
// 实现排序需要的三个方法
// 在其他排序方法中被调用
type Sorter interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

type DebugSorter interface {
	Sorter
	Index(...int) string
}

// sort float64
type float64Sorter []float64
func (s float64Sorter) Len() int           { return len(s) }
func (s float64Sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s float64Sorter) Less(i, j int) bool { return s[i] < s[j] }
func (s float64Sorter) Index(index ...int) string {
	var temp = make(float64Sorter, len(index))
	for i, j := range index {
		temp[i] = s[j]
	}
	return fmt.Sprintf("%+#v", temp)
}

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

func SortSlice(slice interface{}, less func(i, j int) bool, sortFunc ...func(Sorter)) {
	var typ = reflect.TypeOf(slice)
	if typ.Kind() != reflect.Slice {
		return
	}
}
