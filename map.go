package immutable

import (
	"reflect"
	"unsafe"
)

const (
	bucketCount    uint32 = 8
	levels         uint32 = 4
	leafStartCount uint32 = 1
)

// Map is an immutable hash map with copy-on-write semantics.
// Adding to or deleting from the map returns a new map instance.
// Since the map is immutable, it is safe to use from multiple
// concurrent threads without locks or other synchronization.
//
// Copying the map is cheap, but since it is a struct, it is not
// done atomically. To get atomic copying, use channels, atomic.Value
// or similar.
//
// Map is different from Go map and sync.Map since it safe to
// copy and is copied by value.
//
// The zero Map is empty and ready for use.
type Map struct {
	leafCount uint32
	capacity  uint32
	size      uint32
	root      bucket
}

// Set adds an entry to a map and returns the updated map.
func (m Map) Set(key, value interface{}) Map {
	hash := hashValue(key)

	if m.capacity == 0 {
		m.leafCount = leafStartCount
		m.capacity = mapCapacity(m.leafCount)
	} else if m.size*2 >= m.capacity {
		m.leafCount *= 2
		m.capacity *= 2
	}

	b := &m.root

	for level := uint32(0); level < levels; level++ {
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

	newValues := make([]elementList, m.leafCount)

	if uint32(len(b.values)) != m.leafCount {
		for _, list := range b.values {
			for _, element := range list {
				hash := hashValue(element.key)
				for l := uint32(0); l < levels; l++ {
					hash /= bucketCount
				}

				valueIndex := hash % m.leafCount
				newList := newValues[valueIndex]
				newList = append(newList, element)
				newValues[valueIndex] = newList
			}
		}
	} else {
		copy(newValues, b.values)
	}

	b.values = newValues

	valueIndex := hash % m.leafCount
	list := b.values[valueIndex]
	list = append(list[:0:0], list...)

	for i, e := range list {
		if e.key == key {
			e.value = value
			list[i] = e
			b.values[valueIndex] = list
			return m
		}
	}

	list = append(list, element{key, value})
	b.values[valueIndex] = list
	m.size++
	return m
}

// Get retrieves a value from the map.
func (m Map) Get(key interface{}) (interface{}, bool) {
	if m.capacity == 0 {
		return nil, false
	}

	hash := hashValue(key)

	b := &m.root
	for level := uint32(0); level < levels; level++ {
		bucketIndex := hash % bucketCount
		next := b.buckets[bucketIndex]
		if next == nil {
			return nil, false
		}
		b = next
		hash /= bucketCount
	}

	if len(b.values) == 0 {
		return nil, false
	}

	valueIndex := hash % uint32(len(b.values))
	list := b.values[valueIndex]

	for _, e := range list {
		if e.key == key {
			return e.value, true
		}
	}

	return nil, false
}

// Delete returns a map without entries matching the key.
// If no entry matches, the original map is returned.
func (m Map) Delete(key interface{}) Map {
	if m.capacity == 0 {
		return m
	}

	hash := hashValue(key)

	root := m.root
	b := &root

	for level := uint32(0); level < levels; level++ {
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
	newValues := make([]elementList, m.leafCount)
	copy(newValues, b.values)
	b.values = newValues

	valueIndex := hash % uint32(len(b.values))
	list := b.values[valueIndex]
	list = append(elementList{}, list...)

	for i, e := range list {
		if e.key == key {
			list = append(list[0:i], list[i+1:]...)
			b.values[valueIndex] = list
			return Map{
				size: m.size - 1,
				root: root,
			}
		}
	}
	return m
}

// Range calls visitor for each element in the map.
// If visitor returns false, the iteration stops.
// Since the map is immutable, it will not change during iteration.
func (m *Map) Range(visitor func(key, value interface{}) bool) {
	m.root.visit(visitor)
}

// Size returns the number of elements in the map.
func (m *Map) Size() uint32 {
	return m.size
}

func (b *bucket) visit(visitor func(key, value interface{}) bool) bool {
	if len(b.values) > 0 {
		for _, list := range b.values {
			for _, e := range list {
				keepGoing := visitor(e.key, e.value)
				if !keepGoing {
					return false
				}
			}
		}
	} else {
		for _, child := range b.buckets {
			if child == nil {
				continue
			}
			keepGoing := child.visit(visitor)
			if !keepGoing {
				return false
			}
		}
	}
	return true
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

	default:
		t := reflect.TypeOf(key)
		if !t.Comparable() {
			panic("Key must be comparable")
		}

		iface := (*ifaceWords)(unsafe.Pointer(&key))
		ptr := iface.data

		size := t.Size()
		bytes = (*[512]uint8)(ptr)[:size:size]
	}

	return hashFunc(bytes)
}

func mapCapacity(leafCount uint32) uint32 {
	capacity := uint32(1)
	for level := uint32(0); level < levels; level++ {
		capacity *= bucketCount
	}
	capacity *= leafCount
	return capacity
}

// Hack!
// ifaceWords is interface{} internal representation, copied
// from sync.atomic.
type ifaceWords struct {
	_    unsafe.Pointer
	data unsafe.Pointer
}
