package gosort

func HeapSort(s Sorter) {
	registerDebugIndexFunc(s)
	heapSort(s.Len, s.Less, s.Swap)
}

// 二叉堆上浮操作
func floatUp(childIndex int, binaryHeap []int, less func(int, int) bool) {
	// 如果子节点为根节点，跳过
	if childIndex == 0 {
		return
	}

	// 如果子节点不是根节点
	// 则比较子节点数据与父节点数据

	// parentIndex: 在二叉堆中的父节点 index
	var parentIndex = (childIndex - 1) / 2

	// 如果插入的子节点数据小于父节点数据
	// 那么子节点数据和父节点数据交换
	// 直到上浮过程完毕
	if less(binaryHeap[childIndex], binaryHeap[parentIndex]) {
		binaryHeap[parentIndex], binaryHeap[childIndex] = binaryHeap[childIndex], binaryHeap[parentIndex]

		// 交换完毕之后，继续上浮比较
		floatUp(parentIndex, binaryHeap, less)
	}
}

// 二叉堆下沉操作
func sinkDown(parentIndex int, binaryHeap []int, less func(int, int) bool) {
	var l = len(binaryHeap)
	var leftChildIndex = 2*parentIndex + 1

	// 没有子节点了
	if leftChildIndex >= l {
		return
	}

	// 左子节点较小
	if less(binaryHeap[leftChildIndex], binaryHeap[parentIndex]) {
		binaryHeap[parentIndex], binaryHeap[leftChildIndex] = binaryHeap[leftChildIndex], binaryHeap[parentIndex]

		// 交换完毕，继续下沉比较
		sinkDown(leftChildIndex, binaryHeap, less)
	}

	// 左子节点较大
	var rightChildIndex = 2*parentIndex + 2
	// 右子节点不存在
	if rightChildIndex >= l {
		return
	}

	// 右子节点存在
	// 右子节点较小
	if less(binaryHeap[rightChildIndex], binaryHeap[parentIndex]) {
		binaryHeap[parentIndex], binaryHeap[rightChildIndex] = binaryHeap[rightChildIndex], binaryHeap[parentIndex]

		// 交换完毕，继续下沉比较
		sinkDown(rightChildIndex, binaryHeap, less)
	}
}

// 创建二叉堆
func createBinaryHeap(
	length int,
	less func(int, int) bool,
) []int {
	var binaryHeap = make([]int, length)
	// 二叉堆排序
	for i := 0; i < length; i++ {
		// childIndex: 在二叉堆中的子节点 index
		var childIndex = i
		// 子节点先把对应的原始数组 index 数据插入
		binaryHeap[childIndex] = i

		// 上浮调整成最大堆
		floatUp(childIndex, binaryHeap, less)
	}

	return binaryHeap
}

// 二叉堆排序
func heapSort(
	len func() int,
	less func(int, int) bool,
	swap func(int, int),
) {
	var l = len()
	var originIndex = make([]int, l)
	for i := range originIndex {
		originIndex[i] = i
	}

	// 创建二叉堆
	var binaryHeap = createBinaryHeap(l, less)

	for i := 0; i < l; i++ {
		binaryHeap[0], binaryHeap[l-i-1] = binaryHeap[l-i-1], binaryHeap[0]
		sinkDown(0, binaryHeap[:l-i-1], less)
	}

	tidy(binaryHeap, swap)
	reverse(len, swap)
}
