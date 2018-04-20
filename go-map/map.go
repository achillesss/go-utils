package gomap

import (
	"fmt"
	"reflect"
)

type GoMap struct {
	instance          interface{}
	addChan           chan map[interface{}]interface{}
	afterAddChan      chan struct{}
	delChan           chan interface{}
	afterDelChan      chan struct{}
	queryChan         chan interface{}
	queryRespChan     chan map[interface{}]interface{}
	interfaceChan     chan struct{}
	interfaceRespChan chan interface{}
	setChan           chan interface{}
	afterSetChan      chan struct{}
	checkLengthChan   chan struct{}
	lengthChan        chan int
	dropChan          chan struct{}
}

// NewMap creates a map
func NewMap(srcMap interface{}) *GoMap {
	srcType := reflect.TypeOf(srcMap)
	if srcType.Kind() != reflect.Map {
		panic("src not map")
	}
	var m GoMap
	m.instance = srcMap
	m.addChan = make(chan map[interface{}]interface{})
	m.afterAddChan = make(chan struct{})
	m.delChan = make(chan interface{})
	m.afterDelChan = make(chan struct{})
	m.queryChan = make(chan interface{})
	m.queryRespChan = make(chan map[interface{}]interface{})
	m.interfaceChan = make(chan struct{})
	m.interfaceRespChan = make(chan interface{})
	m.setChan = make(chan interface{})
	m.afterSetChan = make(chan struct{})
	m.lengthChan = make(chan int)
	m.checkLengthChan = make(chan struct{})
	m.dropChan = make(chan struct{})
	return &m
}

func isMap(src interface{}) bool {
	return reflect.TypeOf(src).Kind() == reflect.Map
}

// MapHandler handles map
func (gm GoMap) Handler() {
	// mapValue := reflect.ValueOf(gm.instance)
	mapType := reflect.TypeOf(gm.instance)
	mapValue := reflect.ValueOf(gm.instance)
	if mapValue.IsNil() {
		mapValue = reflect.MakeMap(mapType)
	}

	keysType := mapType.Key()
	valueType := mapType.Elem()

	for {
		select {

		// add
		case m := <-gm.addChan:
			for k, v := range m {
				kv := reflect.ValueOf(k)
				vv := reflect.ValueOf(v)
				if vv.IsValid() {
					mapValue.SetMapIndex(kv, vv)
				}
			}
			gm.afterAddChan <- struct{}{}

		// delete
		case m := <-gm.delChan:
			mv := reflect.ValueOf(m)
			mapValue.SetMapIndex(mv, reflect.Value{})
			gm.afterDelChan <- struct{}{}

		// query
		case m := <-gm.queryChan:
			kt := reflect.TypeOf(m)
			if kt.Kind() != keysType.Kind() {
				gm.queryRespChan <- nil
				continue
			}

			kv := reflect.ValueOf(m)
			v := mapValue.MapIndex(kv)
			newV := reflect.New(valueType)

			if !v.IsValid() {
				newV.Elem().Set(reflect.Zero(valueType))
			} else {
				newV.Elem().Set(v)
			}
			gm.queryRespChan <- map[interface{}]interface{}{m: newV.Elem().Interface()}

		// change to interface{}
		case <-gm.interfaceChan:
			gm.interfaceRespChan <- mapValue.Interface()

		// set map
		case m := <-gm.setChan:
			mapValue = reflect.ValueOf(m)
			gm.afterSetChan <- struct{}{}

		case <-gm.checkLengthChan:
			gm.lengthChan <- mapValue.Len()
		// drop, interrupt select loop
		case <-gm.dropChan:
			return
		}

	}
}

// Add add key: value to map
func (gm *GoMap) Add(key, value interface{}) {
	gm.addChan <- map[interface{}]interface{}{key: value}
	<-gm.afterAddChan
}

func (gm *GoMap) Delete(key interface{}) {
	gm.delChan <- key
	<-gm.afterDelChan
}

func (gm *GoMap) pickQueryResp(key interface{}) interface{} {
	for resp := range gm.queryRespChan {
		if resp == nil {
			break
		}
		for k, v := range resp {
			if reflect.DeepEqual(key, k) {
				return v
			}
		}
		gm.queryRespChan <- resp
	}
	return nil
}

func (gm *GoMap) Set(srcMap interface{}) {
	gm.setChan <- srcMap
	<-gm.afterSetChan
}

func (gm *GoMap) Query(key interface{}, dst interface{}) error {
	gm.queryChan <- key
	v := gm.pickQueryResp(key)
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return fmt.Errorf("bad dst type")
	}

	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		panic("dst not pointer")
	}

	dv.Elem().Set(reflect.ValueOf(v))
	return nil
}

func (gm *GoMap) Interface() interface{} {
	gm.interfaceChan <- struct{}{}
	return <-gm.interfaceRespChan
}

func (gm GoMap) Len() int {
	gm.checkLengthChan <- struct{}{}
	return <-gm.lengthChan
}

// Close means no other coming actions after it.
func (gm *GoMap) Close() {
	gm.dropChan <- struct{}{}
}
