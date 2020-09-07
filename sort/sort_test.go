package gosort

import (
	"math/rand"
	"testing"
	"time"

	"github.com/achillesss/go-utils/log"
)

func TestSelectionSort(t *testing.T) {
	var sortMethods = []func(Sorter){
		SelectionSort,
		InsertionSort,
		BubbleSort,
		ShellSort,

		MergeSort,
		QuickSort,
		HeapSort,
	}

	var sliceLength = 20
	rand.Seed(time.Now().UnixNano())
	var randomSlice = make([]float64, sliceLength)
	for i := range randomSlice {
		randomSlice[i] = float64(rand.Int63n(1000))
	}
	log.Infofln("RandomSlice: %+v", randomSlice)

	for i, method := range sortMethods {
		var src = make([]float64, sliceLength)
		copy(src, randomSlice)
		SortFloat64(src, method)
		log.Infofln("medhod[%d]: %+v", i, src)
	}
}
