package gosort

var indexDebug func(...int) string

func registerDebugIndexFunc(s Sorter) {
	indexDebug = s.(DebugSorter).Index
}
