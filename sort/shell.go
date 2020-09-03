package gosort

func shellSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()

	// 逻辑同插入排序
	// 只是增加了插入步长 gap
	var sortFunc = func(insertIndex, gap int) {
		for i := 0; i < insertIndex; i += gap {
			if less(insertIndex, i) {
				swap(insertIndex, i)
			}
		}
	}

	for gap := l >> 1; gap > 0; gap >>= 1 {
		for i := gap; i < l; i++ {
			sortFunc(i, gap)
		}
	}
}

func ShellSort(s Sorter) {
	shellSort(s.Len, s.Less, s.Swap)
}
