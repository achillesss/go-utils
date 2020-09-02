package gosort

// Selection Sort 选择排序
func selectionSort(s Sorter) {
	var l = s.Len()

	// 遍历未排序的部分
	// 找出最小值
	// 放在第一位
	// firstIndex：每次循环时，无序数组的第一个元素，其值等于循环次数
	var sortFunc = func(firstIndex int) {
		var min = firstIndex
		for i := firstIndex + 1; i < l; i++ {
			if s.Less(i, min) {
				min = i
			}
		}
		s.Swap(min, firstIndex)
	}

	for i := 0; i < l; i++ {
		sortFunc(i)
	}
}
