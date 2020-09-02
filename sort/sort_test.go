package gosort
import (
	"fmt"
	"testing"
)

func TestSelectionSort(t *testing.T) {
	var sortMethods = []func(Sorter){
		SelectionSort,
		InsertionSort,
		BubbleSort,
	}

	for i, method := range sortMethods {
		fmt.Printf("method: %d\n", i)
		var src = []float64{4, 3, 2, 1, 200, 300, 100}
		fmt.Printf("before: %+v\n", src)
		SortFloat64(src, method)
		fmt.Printf("after: %+v\n", src)
	}
}
