package gomap

import (
	"fmt"
	"testing"

	"github.com/achillesss/go-utils/log"
)

func TestMapInt(t *testing.T) {
	srcMap := make(map[int]int)
	m := NewMap(srcMap)
	var q int

	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	m.Query(1, &q)
	if q != 0 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}

	m.Add(1, 1)
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	m.Query(1, &q)
	if q != 1 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())

	m.Delete(1)
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())

	m.Query(1, &q)
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	if q != 0 {
		t.Errorf("%s failed.", log.FuncName())
		return
	}

	m.Set(map[int]int{11: 11, 22: 22})
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	result := make(map[int]int)
	m.BatchQuery([]int{11, 22}, &result)
	fmt.Printf("result: %+v\n", result)

	m.BatchAdd(map[int]int{33: 33, 44: 44})
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
}

func TestMapString(t *testing.T) {
	var srcMap map[string]string
	m := NewMap(srcMap)
	var q string

	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	m.Query("1", &q)
	if q != "" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}

	m.Add("1", "1")
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())

	m.Query("1", &q)
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	if q != "1" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())

	m.Delete("1")
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())

	m.Query("1", &q)
	fmt.Printf("map: %+v, length: %d\n", m.Interface(), m.Len())
	if q != "" {
		t.Errorf("%s failed.", log.FuncName())
		return
	}

}

func TestMapValueStruct(t *testing.T) {
	type A struct {
		a bool
		b int
		c string
		d []int
	}
	var srcMap map[int]*A
	var a A
	a.a = true
	a.b = 2
	a.c = "hello"
	a.d = []int{3, 4}

	m := NewMap(srcMap)
	var q *A
	m.Query(9, &q)
	fmt.Printf("q: %+#v\n", q)
	m.Add(9, &a)
	m.Query(9, &q)
	if q == nil || q.a != a.a || q.b != a.b && q.c != a.c || (len(q.d) != len(q.d)) || q.d[0] != q.d[0] || q.d[1] != q.d[1] {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
	m.Delete(0)
	m.Query(0, &q)
	if q != nil {
		t.Errorf("%s failed.", log.FuncName())
		return
	}
}
