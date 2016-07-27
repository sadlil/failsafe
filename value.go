package failsafe

import (
	"reflect"
	"sync"
)

type value struct {
	mu  sync.RWMutex
	val map[string]*abstractValue
}

func newValue() *value {
	return &value{
		val: make(map[string]*abstractValue),
	}
}

func (v *value) Get(key string) (*abstractValue, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	value, ok := v.val[key]
	if !ok {
		return nil, false
	}
	return value, true
}

func (v *value) Put(key string, val interface{}) *value {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.val[key] = newAbstractValue(val)
	return v
}

func (v *value) Remove(key string) *value {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.val, key)
	return v
}

func (v *value) Clear() *value {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.val = make(map[string]*abstractValue)
	return v
}

func (v *value) Len() int {
	return len(v.val)
}

// underlying value types that holds the data
type abstractValue struct {
	value interface{}

	reflectValue reflect.Value
	reflectKind  reflect.Kind
}

func newAbstractValue(v interface{}) *abstractValue {
	return &abstractValue{
		value:        v,
		reflectValue: reflect.ValueOf(v),
		reflectKind:  reflect.ValueOf(v).Kind(),
	}
}

func (a *abstractValue) Interface() interface{} {
	return a.value
}

func (a *abstractValue) ReflectValue() reflect.Value {
	return a.reflectValue
}

func (a *abstractValue) String() string {
	return a.reflectValue.String()
}

func (a *abstractValue) Bool() bool {
	return a.reflectValue.Bool()
}

func (a *abstractValue) Uint() uint {
	return uint(a.reflectValue.Uint())
}

func (a *abstractValue) Uint8() uint8 {
	return uint8(a.reflectValue.Uint())
}

func (a *abstractValue) Uint16() uint16 {
	return uint16(a.reflectValue.Uint())
}

func (a *abstractValue) Uint32() uint32 {
	return uint32(a.reflectValue.Uint())
}

func (a *abstractValue) Uint64() uint64 {
	return a.reflectValue.Uint()
}

func (a *abstractValue) Int() int {
	return int(a.reflectValue.Int())
}

func (a *abstractValue) Int8() int8 {
	return int8(a.reflectValue.Int())
}

func (a *abstractValue) Int16() int16 {
	return int16(a.reflectValue.Int())
}

func (a *abstractValue) Int32() int32 {
	return int32(a.reflectValue.Int())
}

func (a *abstractValue) Int64() int64 {
	return a.reflectValue.Int()
}

func (a *abstractValue) Float32() float32 {
	return float32(a.reflectValue.Float())
}

func (a *abstractValue) Float64() float64 {
	return a.reflectValue.Float()
}

func (a *abstractValue) Uintptr() uintptr {
	return a.reflectValue.Pointer()
}

func (a *abstractValue) Complex64() complex64 {
	return complex64(a.reflectValue.Complex())
}

func (a *abstractValue) Complex128() complex128 {
	return a.reflectValue.Complex()
}

func (a *abstractValue) Bytes() []byte {
	return a.reflectValue.Bytes()
}

func (a *abstractValue) Slice() []*abstractValue {
	if a.reflectKind != reflect.Slice {
		panic("Slice() given a non-slice type")
	}
	ret := make([]*abstractValue, a.reflectValue.Len())
	for i := 0; i < a.reflectValue.Len(); i++ {
		ret[i] = newAbstractValue(a.reflectValue.Index(i).Interface())
	}
	return ret
}
