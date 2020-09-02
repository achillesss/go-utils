package gosort

// Insertion Sort 插入排序
func insertionSort(s Sorter) {
	var l = s.Len()

	// 要插入的 index
	// 与已经排好序的数组中每一个元素做比较
	// 如果更小，就往前排
	// insertIndex: 每次循环需要插入的元素，其值等于循环次数
	var sortFunc = func(insertIndex int) {
		for i := 0; i < insertIndex; i++ {
			if s.Less(insertIndex, i) {
				s.Swap(insertIndex, i)
			}
		}
	}

	for i := 0; i < l; i++ {
		sortFunc(i)
	}
}
