package quick

import (
	"sort"

	"github.com/achillesss/go-utils/arithmetic/sort/tidy"
)

func Sort(s sort.Interface) {
	quickSort(s.Len, s.Less, s.Swap)
}

func quickSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()
	if l == 1 {
		return
	}

	// 原始数组的 index
	var src = make([]int, l)

	// 初始化
	for i := range src {
		src[i] = i
	}

	// 对原始数组的 index 排序
	// 返回一个已经排序好的 index 数组
	var sortedIndexArr = quickSortTemp(src, less)

	// 根据已排序的 index 整理原始数组
	tidy.Tidy(sortedIndexArr, swap)
}

func quickSortTemp(src []int, less func(int, int) bool) []int {
	var l = len(src)
	if l < 2 {
		return src
	}

	var temp = make([]int, l)
	var n int = 0
	var m int = l - 1

	var pivotIndex = l - 1
	for i := 0; i < l-1; i++ {
		if less(src[i], src[pivotIndex]) {
			temp[n] = src[i]
			n++
			continue
		}
		temp[m] = src[i]
		m--
	}

	temp[n] = src[pivotIndex]
	if len(temp) < 3 {
		return temp
	}

	var merged = append(quickSortTemp(temp[:n], less), quickSortTemp(temp[n:], less)...)
	return merged
}
