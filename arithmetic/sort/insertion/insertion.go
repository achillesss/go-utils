package insertion

import (
	"sort"
)

func Sort(s sort.Interface) {
	insertionSort(s.Len, s.Less, s.Swap)
}

// Insertion Sort 插入排序
func insertionSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()

	// 要插入的元素
	// 与已经排好序的数组中每一个元素做比较
	// 如果更小，就往前排
	// insertIndex: 每次循环需要插入的元素，其值等于循环次数
	var sortFunc = func(insertIndex int) {
		for i := 0; i < insertIndex; i++ {
			if less(insertIndex, i) {
				swap(insertIndex, i)
			}
		}
	}

	for i := 1; i < l; i++ {
		sortFunc(i)
	}
}
