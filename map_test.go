package immutable

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetEmptyByString(t *testing.T) {
	var m Map
	v, ok := m.Get("something")
	if v != nil || ok {
		t.Fail()
	}
}

func TestGetEmptyByInteger(t *testing.T) {
	var m Map
	v, ok := m.Get(42)
	if v != nil || ok {
		t.Fail()
	}
}

func TestSetGetByString(t *testing.T) {
	var m Map
	key := "kawonka"
	m = m.Set(key, 124)
	v, ok := m.Get(key)
	if v == nil || !ok {
		t.Fail()
	}
	if v.(int) != 124 {
		t.Fail()
	}
}

func TestSetGetByInt(t *testing.T) {
	var m Map
	key := 987
	m = m.Set(key, 4711)
	v, ok := m.Get(key)
	if v == nil || !ok {
		t.Fail()
	}
	if v.(int) != 4711 {
		t.Fail()
	}
}

func TestSetByArrayFails(t *testing.T) {
	var m Map
	key := [12]float32{}
	defer func() {
		if recover() == nil {
			t.Fail()
		}
	}()
	m = m.Set(key, 4711)
}

func TestResetSameKey(t *testing.T) {
	var m Map
	m = m.Set("hej", 2)
	m = m.Set("hej", 3)
	fetched, _ := m.Get("hej")
	if fetched != 3 {
		t.Fail()
	}
}

func TestNoSharing(t *testing.T) {
	var a Map
	a = a.Set("hej", "svej")
	b := a.Set("hej", "hoj")
	b = b.Set("vem", "där")
	aHej, _ := a.Get("hej")
	bHej, _ := b.Get("hej")
	if aHej == bHej {
		t.Fail()
	}
	_, ok := a.Get("vem")
	if ok {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	var m Map
	m = m.Set(9876, 1234)
	m = m.Set("number", 42)
	_, ok := m.Get("number")
	if !ok {
		t.Fail()
	}
	d := m.Delete("number")
	_, ok = d.Get("number")
	if ok {
		t.Fail()
	}
	_, ok = m.Get("number")
	if !ok {
		t.Fail()
	}
}

func TestAddMany(t *testing.T) {
	var m Map
	for i := 0; i < 1000; i++ {
		m = m.Set(i, i)
	}
	for i := 0; i < 1000; i++ {
		v, ok := m.Get(i)
		if !ok || v.(int) != i {
			t.Fail()
		}
	}
	for i := 1000; i < 2000; i++ {
		_, ok := m.Get(i)
		if ok {
			t.Fail()
		}
	}
}

func BenchmarkAddIntsImmutableMap(b *testing.B) {
	var m atomic.Value
	m.Store(Map{})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				old := m.Load().(Map)
				updated := old.Set(num, num)
				m.Store(updated)
			}
		}
	})
}

func BenchmarkAddIntsGoMap(b *testing.B) {
	m := map[int]int{}
	var mutex sync.Mutex

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				mutex.Lock()
				m[num] = num
				mutex.Unlock()
			}
		}
	})
}

func BenchmarkAddIntsSyncMap(b *testing.B) {
	m := sync.Map{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				m.Store(num, num)
			}
		}
	})
}

func BenchmarkGetIntsImmutableMap(b *testing.B) {
	var m Map
	for i := 0; i < b.N; i++ {
		num := i % 1024
		m = m.Set(num, num)
	}

	var mm atomic.Value
	mm.Store(m)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				loaded := mm.Load().(Map)
				_, _ = loaded.Get(num)
			}
		}
	})
}

func BenchmarkGetIntsGoMap(b *testing.B) {
	m := map[int]int{}
	for i := 0; i < b.N; i++ {
		num := i % 1024
		m[num] = num
	}

	var mutex sync.Mutex

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				mutex.Lock()
				_, _ = m[num]
				mutex.Unlock()
			}
		}
	})
}

func BenchmarkGetIntsSyncMap(b *testing.B) {
	m := sync.Map{}
	for i := 0; i < b.N; i++ {
		num := i % 1024
		m.Store(num, num)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				num := i % 1024
				_, _ = m.Load(num)
			}
		}
	})
}
