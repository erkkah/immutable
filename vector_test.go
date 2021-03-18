package immutable

import "testing"

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
