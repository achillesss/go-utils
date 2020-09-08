package gosort

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/achillesss/go-utils/arithmetic/sort/bubble"
	"github.com/achillesss/go-utils/arithmetic/sort/heap"
	"github.com/achillesss/go-utils/arithmetic/sort/insertion"
	"github.com/achillesss/go-utils/arithmetic/sort/merge"
	"github.com/achillesss/go-utils/arithmetic/sort/quick"
	"github.com/achillesss/go-utils/arithmetic/sort/selection"
	"github.com/achillesss/go-utils/arithmetic/sort/shell"
	"github.com/achillesss/go-utils/log"
)

func TestSelectionSort(t *testing.T) {
	var sortMethods = []func(sort.Interface){
		selection.Sort,
		insertion.Sort,
		bubble.Sort,
		shell.Sort,

		merge.Sort,
		quick.Sort,
		heap.Sort,
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
