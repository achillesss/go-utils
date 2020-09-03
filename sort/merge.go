package gosort

// 递归版
func MergeSortRecursive(s Sorter) {
	mergeSortRecursive(s.Len, s.Less, s.Swap)
}

func mergeSortRecursive(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()
	if l == 1 {
		return
	}

	var src = make([]int, l)

	// 老的 index 上经过 swap 之后
	// 对应现在的 index
	var indexRealPlace = make(map[int]int, l)

	for i := range src {
		src[i] = i
		indexRealPlace[i] = i
	}

	var sortedIndexArr = mergeSortRecursiveTemp(src, less)
	for i, srcIndex := range sortedIndexArr {
		// 老的 index 对应的实际 index
		var srcIndexRealPlace = indexRealPlace[srcIndex]

		// 上面的 index
		// 要与此 index 交换位置
		var realIndex = i

		// 交换
		swap(srcIndexRealPlace, realIndex)

		// 交换位置之后
		// 要更新交换记录
		// 1. 老的 index 所在位置被更新成 realIndex
		// 2. 与老 index 交换位置的 realIndex 被老 index 实际上对应的 index 替换，即
		// 被 srcIndexRealPlace 替换
		// 原本 srcIndex: srcIndexRealPlace, realIndex: realIndex
		// 交换之后 srcIndex: realIndex, realIndex: srcIndexRealPlace
		indexRealPlace[srcIndex], indexRealPlace[realIndex] = realIndex, srcIndexRealPlace
	}
}

// 归并排序之后
// 获得一个被排序好的 index 列表
// 外面需根据此列表，使用 Swap 方法来对原始列表排序
func mergeSortRecursiveTemp(srcArr []int, less func(int, int) bool) []int {
	var l = len(srcArr)
	if l < 2 {
		return srcArr
	}

	// mid_length
	var ml = (l-1)/2 + 1

	// Divide
	var arrI = mergeSortRecursiveTemp(srcArr[:ml], less)
	var arrJ = mergeSortRecursiveTemp(srcArr[ml:], less)

	// Merge
	return mergeSortRecursiveMerge(arrI, arrJ, less)
}

func mergeSortRecursiveMerge(arrI, arrJ []int, less func(int, int) bool) []int {
	if arrI == nil {
		return arrJ
	}

	if arrJ == nil {
		return arrJ
	}

	var (
		i    int
		j    int
		lI   = len(arrI)
		lJ   = len(arrJ)
		temp = make([]int, 0, lI+lJ)
	)

	for {
		if i >= lI && j >= lJ {
			break
		}

		if i >= lI {
			temp = append(temp, arrJ[j:]...)
			break
		}

		if j >= lJ {
			temp = append(temp, arrI[i:]...)
			break
		}

		if less(arrI[i], arrJ[j]) {
			temp = append(temp, arrI[i])
			i++
			continue
		}

		temp = append(temp, arrJ[j])
		j++
	}

	return temp
}

// 迭代版
func mergeSortIterative() {

}
