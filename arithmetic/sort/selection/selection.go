package selection

import "sort"

func Sort(s sort.Interface) {
	selectionSort(s.Len, s.Less, s.Swap)
}

// Selection Sort 选择排序
func selectionSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()

	// 遍历未排序的部分
	// 找出最小值
	// 放在第一位
	// firstIndex：每次循环时，无序数组的第一个元素，其值等于循环次数
	var sortFunc = func(firstIndex int) {
		var min = firstIndex
		for i := firstIndex + 1; i < l; i++ {
			if less(i, min) {
				min = i
			}
		}
		swap(min, firstIndex)
	}

	for i := 0; i < l; i++ {
		sortFunc(i)
	}
}
