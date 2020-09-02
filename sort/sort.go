package gosort

type Sorter interface {
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

// sort float64
type float64Sorter []float64
func (s float64Sorter) Len() int           { return len(s) }
func (s float64Sorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s float64Sorter) Less(i, j int) bool { return s[i] < s[j] }

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