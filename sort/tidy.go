package gosort

// 根据已经排序完的 index 列表
// 使用原数组的 swap 方法对原数组排序
func tidy(sortedIndexArr []int, swap func(int, int)) {
	var l = len(sortedIndexArr)

	// 记录每一个原始数组的 index 当前对应真实的原始数组 index
	var srcIndexMap = make(map[int]int, l)

	// 记录被交换之后的原始 index 顺序
	var swappedIndexMap = make(map[int]int, l)

	// 初始化
	// 原始数字没有调用 swap 方法
	// 所以原始数组的 index 所在的位置即 index 本身
	for i := range sortedIndexArr {
		srcIndexMap[i] = i
		swappedIndexMap[i] = i
	}

	// swappedPlace: 被交换到已排好序的数组的何处
	// srcIndex: 原始数组的 index
	for swappedPlace, srcIndex := range sortedIndexArr {
		// 找到原始数组 index 所在的真实 index 位置
		var srcIndexPlace = srcIndexMap[srcIndex]
		var realSrcIndex = swappedIndexMap[swappedPlace]

		if srcIndexPlace == swappedPlace {
			continue
		}

		swap(srcIndexPlace, swappedPlace)

		// 交换位置之后
		// 要更新交换记录
		// 更新原数组 index 所在位置
		srcIndexMap[srcIndex] = swappedPlace
		srcIndexMap[realSrcIndex] = srcIndexPlace

		// 更新排序数组中 index 对应的原数组 index
		swappedIndexMap[swappedPlace] = srcIndex
		swappedIndexMap[srcIndexPlace] = realSrcIndex
	}
}

func reverse(len func() int, swap func(int, int)) {
	var l = len()
	for i := 0; i <= l/2-1; i++ {
		swap(i, l-i-1)
	}
}
