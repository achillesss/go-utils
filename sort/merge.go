package gosort

import (
	"flag"
)
var MergeSortVersion = flag.Int(
	"mergeVersion",
	1,
	`归并排序使用的排序版本。
1. 递归版(默认)
2. 迭代版`,
)

func MergeSort(s Sorter) {
	mergeSort(s.Len, s.Less, s.Swap)
}

func mergeSort(len func() int, less func(int, int) bool, swap func(int, int)) {
	var l = len()
	if l == 1 {
		return
	}

	// 原始数组的 index
	var src = make([]int, l)

	// 老的 index 会被 swap 若干次
	// 所以需要一个变量来记住每一个老的 index 与哪一个 index 交换了位置
	// 初始情况下，老的index对应的位置就是它本身
	var indexRealPlace = make(map[int]int, l)

	// 初始化
	for i := range src {
		src[i] = i
		indexRealPlace[i] = i
	}

	// 对原始数组的 index 排序
	// 返回一个已经排序好的 index 数组

	var sortedIndexArr []int
	switch *MergeSortVersion {
	case 2:
		// 迭代版
		sortedIndexArr = mergeSortIterativeTemp(src, less)
	default:
		// 递归版
		sortedIndexArr = mergeSortRecursiveTemp(src, less)
	}

	// 根据已经排序完的 index 列表
	// 使用原数组的 swap 方法对原数组排序
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
func mergeSortMerge(arrI, arrJ []int, less func(int, int) bool) []int {
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

// 递归版
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
	return mergeSortMerge(arrI, arrJ, less)
}

// 迭代版
func mergeSortIterativeTemp(src []int, less func(int, int) bool) []int {
	var l = len(src)
	if l < 2 {
		return src
	}

	var groupSize = 1
	for groupSize < l {
		var sortedArr = make([]int, 0, l)
		for i := 0; i < l; i += 2 * groupSize {
			var arrI []int
			var arrJ []int

			if i+groupSize >= l {
				arrI = src[i:]
			} else {
				arrI = src[i : i+groupSize]
				if i+groupSize*2 >= l {
					arrJ = src[i+groupSize:]
				} else {
					arrJ = src[i+groupSize : i+groupSize*2]
				}
			}

			sortedArr = append(sortedArr, mergeSortMerge(arrI, arrJ, less)...)
		}

		src = sortedArr
		groupSize <<= 1
	}

	return src
}
