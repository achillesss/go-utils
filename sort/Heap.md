# 堆排序（Heap Sort)

[来自 Wiki 的排序介绍](https://zh.wikipedia.org/wiki/%E5%A0%86%E6%8E%92%E5%BA%8F)

> 在堆的[数据结构](https://zh.wikipedia.org/wiki/資料結構)中，堆中的最大值总是位于根节点（在优先队列中使用堆的话堆中的最小值位于根节点）。堆中定义以下几种操作：
>
> - 最大堆调整（Max Heapify）：将堆的末端子节点作调整，使得子节点永远小于父节点
> - 创建最大堆（Build Max Heap）：将堆中的所有数据重新排序
> - 堆排序（HeapSort）：移除位在第一个数据的根节点，并做最大堆调整的[递归](https://zh.wikipedia.org/wiki/遞迴)运算

时间复杂度：O(nlogn)

空间复杂度：O(1)

