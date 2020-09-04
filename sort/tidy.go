package gosort

// 根据已经排序完的 index 列表
// 使用原数组的 swap 方法对原数组排序
func tidy(sortedIndex []int, swap func(int, int)) {
	var l = len(sortedIndex)

	// 老的 index 会被 swap 若干次
	// 所以需要一个变量来记住每一个老的 index 与哪一个 index 交换了位置
	// 初始情况下，老的index对应的位置就是它本身
	var indexRealPlace = make(map[int]int, l)

	// 初始化
	for i := range sortedIndex {
		indexRealPlace[i] = i
	}

	for i, srcIndex := range sortedIndex {
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

func reverse(len func() int, swap func(int, int)) {
	var l = len()
	for i := 0; i <= l/2-1; i++ {
		swap(i, l-i-1)
	}
}
