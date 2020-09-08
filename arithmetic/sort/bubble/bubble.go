package bubble

import "sort"

func Sort(s sort.Interface) {
	bubbleSort(s.Len, s.Less, s.Swap)
}

// Bubble Sort 冒泡排序
func bubbleSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()

	// 每次循环
	// 都对无序数组[0:len-cnt]
	// 按照冒泡规则
	// 找出一个最大值
	// 放到数组最后
	// loopCnt 为循环次数，也指确定了多少个最大的数
	var sortFunc = func(loopCnt int) {
		for i := 0; i < l-loopCnt-1; i++ {
			if less(i, i+1) {
				continue
			}
			swap(i, i+1)
		}
	}

	for i := 0; i < l; i++ {
		sortFunc(i)
	}
}
