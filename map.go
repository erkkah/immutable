package immutable

import (
	"reflect"
	"unsafe"
)

const (
	bucketCount = 8
	levels      = 4
	leaveCount  = 8
)

// Map is an immutable hash map with copy-on-write semantics
type Map struct {
	root bucket
}

// Set adds an entry to a map and returns the updated map
func (m Map) Set(key, value interface{}) Map {
	hash := hashValue(key)

	root := m.root
	b := &root

	for level := 0; level < levels; level++ {
		bucketIndex := hash % bucketCount

		next := b.buckets[bucketIndex]
		if next == nil {
			next = &bucket{}
		} else {
			next = &bucket{
				next.buckets,
				next.values,
			}
		}
		b.buckets[bucketIndex] = next

		hash /= bucketCount
		b = next
	}

	newValues := make([]elementList, leaveCount, leaveCount)
	copy(newValues, b.values)
	b.values = newValues

	valueIndex := hash % leaveCount
	list := b.values[valueIndex]
	list = append(list[:0:0], list...)

	for i, e := range list {
		if e.key == key {
			e.value = value
			list[i] = e
			b.values[valueIndex] = list
			return Map{root}
		}
	}

	list = append(list, element{key, value})
	b.values[valueIndex] = list
	return Map{root}
}

// Get retrieves a value from the map
func (m Map) Get(key interface{}) (interface{}, bool) {
	hash := hashValue(key)

	b := &m.root
	for level := 0; level < levels; level++ {
		bucketIndex := hash % bucketCount
		next := b.buckets[bucketIndex]
		if next == nil {
			return nil, false
		}
		b = next
		hash /= bucketCount
	}
	valueIndex := hash % leaveCount
	if len(b.values) == 0 {
		return nil, false
	}
	list := b.values[valueIndex]

	for _, e := range list {
		if e.key == key {
			return e.value, true
		}
	}

	return nil, false
}

// Delete returns a map without entries matching the key
func (m Map) Delete(key interface{}) Map {
	hash := hashValue(key)

	root := m.root
	b := &root

	for level := 0; level < levels; level++ {
		bucketIndex := hash % bucketCount

		next := b.buckets[bucketIndex]
		if next == nil {
			return m
		}
		next = &bucket{
			next.buckets,
			next.values,
		}
		b.buckets[bucketIndex] = next

		hash /= bucketCount
		b = next
	}

	if len(b.values) == 0 {
		return m
	}
	newValues := make([]elementList, leaveCount, leaveCount)
	copy(newValues, b.values)
	b.values = newValues

	valueIndex := hash % leaveCount
	list := b.values[valueIndex]
	list = append(elementList{}, list...)

	for i, e := range list {
		if e.key == key {
			list = append(list[0:i], list[i+1:]...)
			b.values[valueIndex] = list
			return Map{root}
		}
	}
	return m
}

type bucket struct {
	buckets [bucketCount]*bucket
	values  []elementList
}

type elementList []element

type element struct {
	key   interface{}
	value interface{}
}

func hashValue(key interface{}) uint32 {
	var bytes []uint8

	// Special fast cases here
	switch val := key.(type) {
	case string:
		bytes = []byte(val)
	case int:
		ptr := unsafe.Pointer(&val)
		const size = unsafe.Sizeof(val)
		bytes = (*[size]uint8)(ptr)[:size:size]
	case int32:
		ptr := unsafe.Pointer(&val)
		const size = unsafe.Sizeof(val)
		bytes = (*[size]uint8)(ptr)[:size:size]
	case int64:
		ptr := unsafe.Pointer(&val)
		const size = unsafe.Sizeof(val)
		bytes = (*[size]uint8)(ptr)[:size:size]
	case float32:
		ptr := unsafe.Pointer(&val)
		const size = unsafe.Sizeof(val)
		bytes = (*[size]uint8)(ptr)[:size:size]
	case float64:
		ptr := unsafe.Pointer(&val)
		const size = unsafe.Sizeof(val)
		bytes = (*[size]uint8)(ptr)[:size:size]
	case []byte:
		bytes = val
	}

	// Generic slow case use reflection
	if len(bytes) == 0 {
		v := reflect.ValueOf(key)
		var ptr unsafe.Pointer = unsafe.Pointer(&key)

		switch v.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Struct:
			fallthrough
		case reflect.Chan:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Interface:
			fallthrough
		case reflect.Map:
			fallthrough
		case reflect.Ptr:
			fallthrough
		case reflect.Slice:
			fallthrough
		case reflect.UnsafePointer:
			panic("Invalid key type")
		}
		t := v.Type()
		if !t.Comparable() {
			panic("Key must be comparable")
		}
		size := t.Size()
		bytes = (*[512]uint8)(ptr)[:size:size]
	}

	var hash uint32

	for _, byte := range bytes {
		hash = hash*31 + uint32(byte)
	}

	return hash
}
