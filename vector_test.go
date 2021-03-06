package immutable

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestGetEmptyVectorSize(t *testing.T) {
	var v Vector
	size := v.Size()
	if size != 0 {
		t.Fail()
	}
}

func TestResizeEmptyVector(t *testing.T) {
	var v Vector
	size := v.Size()
	if size != 0 {
		t.Fail()
	}
	resized := v.Resize(123)
	newSize := resized.Size()
	if newSize != 123 {
		t.Fail()
	}
}

func TestSetValidFirstLevelIndex(t *testing.T) {
	var v Vector
	v = v.Resize(200)
	v = v.Set(3, "3")
	if v.Size() != 200 {
		t.Fail()
	}
}

func TestSetValidSecondLevelIndex(t *testing.T) {
	var v Vector
	v = v.Resize(200)
	v = v.Set(bucketSize, "test")
	if v.Size() != 200 {
		t.Fail()
	}
}

func TestSetGetValidFirstLevelIndex(t *testing.T) {
	var v Vector
	v = v.Resize(200)
	setV := v.Set(3, "3")
	val := setV.Get(3)
	if val != "3" {
		t.Fail()
	}
	val = v.Get(3)
	if val != nil {
		t.Fail()
	}
}

func TestSetGetValidSecondLevelIndex(t *testing.T) {
	var v Vector
	v = v.Resize(200)
	setV := v.Set(bucketSize, "test")
	val := setV.Get(bucketSize)
	if val != "test" {
		t.Fail()
	}
	val = v.Get(bucketSize)
	if val != nil {
		t.Fail()
	}
}

func TestSetOutOfBounds(t *testing.T) {
	var v Vector
	defer func() {
		// Expected failure
		recover()
	}()
	v.Set(2, 2)
	t.Fail()
}

func TestGetOutOfBounds(t *testing.T) {
	var v Vector
	defer func() {
		// Expected failure
		recover()
	}()
	v.Get(2)
	t.Fail()
}

func TestAppendEmpty(t *testing.T) {
	var v Vector
	expected := 4711
	appended := v.Append(expected)
	if appended.Size() != 1 {
		t.Fail()
	}
	val := appended.Get(0)
	if val != expected {
		t.Fail()
	}
}

func TestAppendNonEmpty(t *testing.T) {
	var v Vector
	expected := 4711
	expectedSize := bucketSize*2 + 1
	appended := v.Resize(expectedSize - 1).Append(expected)
	if appended.Size() != expectedSize {
		t.Fail()
	}
	val := appended.Get(expectedSize - 1)
	if val != expected {
		t.Fail()
	}
}

func TestResizeGrowNonEmpty(t *testing.T) {
	v := Vector{}.Resize(bucketSize)
	for i := uint32(0); i < bucketSize; i++ {
		v = v.Set(i, i)
	}
	v = v.Resize(200)
	for i := uint32(0); i < bucketSize; i++ {
		val := v.Get(i)
		if val != i {
			t.Fail()
		}
	}
}

func TestResizeShrinkNonEmpty(t *testing.T) {
	v := Vector{}.Resize(bucketSize * 12)
	for i := uint32(0); i < bucketSize; i++ {
		v = v.Set(i, i)
	}
	v = v.Resize(bucketSize)
	for i := uint32(0); i < bucketSize; i++ {
		val := v.Get(i)
		if val != i {
			t.Fail()
		}
	}
}

func TestSetGetRange(t *testing.T) {
	var v Vector
	var expected [511]int

	v = v.Resize(uint32(len(expected)))
	for i := range expected {
		expected[i] = i
		v = v.Set(uint32(i), i)
		val := v.Get(uint32(i))
		if val != i {
			t.Fail()
		}
	}

	for i := range expected {
		val := v.Get(uint32(i))
		if val != expected[i] {
			t.Fail()
		}
	}
}

func TestSliceValidRange_A(t *testing.T) {
	var v Vector
	var expected [512]int

	v = v.Resize(uint32(len(expected)))
	for i := range expected {
		expected[i] = i
		v = v.Set(uint32(i), i)
	}

	sliced := v.Slice(0, 20)
	expectedSlice := expected[:20]
	for i := range expectedSlice {
		val := sliced.Get(uint32(i))
		if val != expectedSlice[i] {
			t.Fail()
		}
	}
}

func TestSliceValidRange_B(t *testing.T) {
	var v Vector
	var expected [512]int

	v = v.Resize(uint32(len(expected)))
	for i := range expected {
		expected[i] = i
		v = v.Set(uint32(i), i)
	}

	sliced := v.Slice(112, 139)
	expectedSlice := expected[112:139]
	for i := range expectedSlice {
		val := sliced.Get(uint32(i))
		if val != expectedSlice[i] {
			t.Fail()
		}
	}
}

func TestSliceEmptyRanges(t *testing.T) {
	var v Vector
	var expected [512]int

	v = v.Resize(uint32(len(expected)))
	for i := range expected {
		expected[i] = i
		v = v.Set(uint32(i), i)
	}

	sliced := v.Slice(99, 99)
	if sliced.Size() != 0 {
		t.Fail()
	}

	sliced = v.Slice(0, 0)
	if sliced.Size() != 0 {
		t.Fail()
	}

	sliced = v.Slice(1000, 2000)
	if sliced.Size() != 0 {
		t.Fail()
	}
}

func TestSliceLimitedRanges(t *testing.T) {
	var v Vector
	var expected [512]int

	v = v.Resize(uint32(len(expected)))
	for i := range expected {
		expected[i] = i
		v = v.Set(uint32(i), i)
	}

	sliced := v.Slice(511, 701)
	if sliced.Size() != 1 {
		t.Fail()
	}
	if sliced.Get(0) != 511 {
		t.Fail()
	}
}

func TestRangeEmptyVector(t *testing.T) {
	var v Vector

	r := v.Elements()
	for r.Next() {
		t.Fail()
	}
}

func TestVectorRange(t *testing.T) {
	v := Vector{}.Resize(112)
	v = v.Set(42, 42)
	r := v.Elements()
	times := uint32(0)
	for r.Next() {
		if times == 42 {
			if r.Get() != 42 {
				t.Fail()
			}
		}
		times++
	}
	if times != v.Size() {
		t.Fail()
	}
}

const (
	numValues = 1024
)

func BenchmarkRandomAccessFillImmutableVector(b *testing.B) {
	var v atomic.Value
	v.Store(Vector{}.Resize(numValues))
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					old := v.Load().(Vector)
					updated := old.Set(uint32(j), j)
					v.Store(updated)
				}
			}
		}
	})
}

func BenchmarkRandomAccessFillImmutableOwnVector(b *testing.B) {
	v := Vector{}.Resize(numValues)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					v = v.Set(uint32(j), j)
				}
			}
		}
	})
}

func BenchmarkRandomAccessFillGoArray(b *testing.B) {
	var a [numValues]int
	var mutex sync.Mutex
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					mutex.Lock()
					a[j] = j
					mutex.Unlock()
				}
			}
		}
	})
}

func BenchmarkRandomAccessFillCopyGoArray(b *testing.B) {
	var a [numValues]int
	var mutex sync.Mutex
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					mutex.Lock()
					var old = a
					old[j] = j
					a = old
					mutex.Unlock()
				}
			}
		}
	})
}

func BenchmarkRandomAccessGetImmutableVector(b *testing.B) {
	v := Vector{}.Resize(numValues)
	for j := 0; j < numValues; j++ {
		v = v.Set(uint32(j), j)
	}
	var vv atomic.Value
	vv.Store(v)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					vec := vv.Load().(Vector)
					val := vec.Get(uint32(j))
					_ = val
				}
			}
		}
	})
}

func BenchmarkRandomAccessGetImmutableOwnVector(b *testing.B) {
	v := Vector{}.Resize(numValues)
	for j := 0; j < numValues; j++ {
		v = v.Set(uint32(j), j)
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					val := v.Get(uint32(j))
					_ = val
				}
			}
		}
	})
}

func BenchmarkRandomAccessGetGoArray(b *testing.B) {
	var a [numValues]int
	var mutex sync.Mutex
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					mutex.Lock()
					val := a[j]
					_ = val
					mutex.Unlock()
				}
			}
		}
	})
}

func BenchmarkAppendImmutableVector(b *testing.B) {
	var v atomic.Value
	v.Store(Vector{})
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					old := v.Load().(Vector)
					if j == 0 {
						old = old.Resize(0)
					}
					updated := old.Append(j)
					v.Store(updated)
				}
			}
		}
	})
}

func BenchmarkAppendGoSlice(b *testing.B) {
	a := []int{}
	var mutex sync.Mutex
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				for j := 0; j < numValues; j++ {
					mutex.Lock()
					if j == 0 {
						a = []int{}
					}
					a = append(a, j)
					mutex.Unlock()
				}
			}
		}
	})
}

func BenchmarkResizeGoSlice(b *testing.B) {
	a := make([]int, 1000)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				b := make([]int, 10000)
				copy(b, a)
			}
		}
	})
}

func BenchmarkResizeImmutableVector(b *testing.B) {
	a := Vector{}.Resize(1000)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				b := a.Resize(10000)
				_ = b
			}
		}
	})
}
