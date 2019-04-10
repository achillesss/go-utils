package gomap

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type GoMap struct {
	lock         sync.Mutex
	mapType      reflect.Type
	mapValue     reflect.Value
	mapKeyType   reflect.Type
	mapValueType reflect.Type
}

func isMap(src interface{}) bool {
	return reflect.TypeOf(src).Kind() == reflect.Map
}

func NewMap(srcMap interface{}) *GoMap {
	if !isMap(srcMap) {
		panic("src not map")
	}
	var m GoMap
	m.mapType = reflect.TypeOf(srcMap)
	m.mapValue = reflect.ValueOf(srcMap)
	if m.mapValue.IsNil() {
		m.mapValue = reflect.MakeMap(m.mapType)
	}
	m.mapKeyType = m.mapType.Key()
	m.mapValueType = m.mapType.Elem()
	return &m
}

// Add add key: value to map
func (m *GoMap) Add(key, value interface{}) {
	if m == nil {
		return
	}
	m.lock.Lock()
	defer m.lock.Unlock()

	kv := reflect.ValueOf(key)
	vv := reflect.ValueOf(value)
	if vv.IsValid() {
		m.mapValue.SetMapIndex(kv, vv)
	}
}

func (m *GoMap) Delete(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	kv := reflect.ValueOf(key)
	m.mapValue.SetMapIndex(kv, reflect.Value{})
}

func (m *GoMap) Set(srcMap interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.mapValue = reflect.ValueOf(srcMap)
	m.mapType = reflect.TypeOf(srcMap)
	if m.mapValue.IsNil() {
		m.mapValue = reflect.MakeMap(m.mapType)
	}
	m.mapKeyType = m.mapType.Key()
	m.mapValueType = m.mapType.Elem()
}

func (m *GoMap) Query(key interface{}, dst interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	kt := reflect.TypeOf(key)
	if kt.Kind() != m.mapKeyType.Kind() {
		return errors.New("invalid key")
	}

	result := reflect.Zero(m.mapValueType).Interface()
	kv := reflect.ValueOf(key)
	v := m.mapValue.MapIndex(kv)
	if v.IsValid() {
		result = v.Interface()
	}

	dt := reflect.TypeOf(dst)
	if dt.Kind() != reflect.Ptr {
		return fmt.Errorf("dst not a pointer")
	}

	dv := reflect.ValueOf(dst)
	if dv.Kind() != reflect.Ptr {
		panic("dst not a pointer")
	}

	dv.Elem().Set(reflect.ValueOf(result))
	return nil
}

func (m *GoMap) Interface() interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.mapValue.Interface()
}

func (m *GoMap) Len() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.mapValue.Len()
}
