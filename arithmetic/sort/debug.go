package gosort

import "sort"

var indexDebug func(...int) string

func registerDebugIndexFunc(s sort.Interface) {
	indexDebug = s.(DebugSorter).Index
}
