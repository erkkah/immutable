package immutable

// Vector is an immutable vector with copy-on-write semantics.
// Modifying the vector returns a new vector instance.
// Since the vector is immutable, it is safe to use from
// multiple concurrent threads without locks or other
// synchronization.
//
// Copying the vector is cheap, but since it is a struct,
// it is not done atomically. To get atomic copying, use
// atomic.Value or similar.
//
// The zero Vector is empty and ready for use.
type Vector struct {
	// number of adressable elements
	size uint32
	// capacity of storage structure, always >= size
	capacity uint32
	// capacity = bucketSize^(depth), except for zero Vector
	depth uint32
	// offset into storage structure, for slicing
	offset uint32
	// storage root
	root *vectorNode
}

type vectorNode struct {
	values   [bucketSize]interface{}
	children [bucketSize]*vectorNode
}

const (
	bucketBits uint32 = 5
	bucketSize uint32 = 1 << bucketBits
	bucketMask uint32 = bucketSize - 1
)

func (v Vector) Set(index uint32, value interface{}) Vector {
	if index >= v.size {
		panic("Out of bounds vector access")
	}

	index += v.offset

	node := v.root
	nodeIndex := index

	newRoot := &vectorNode{}
	newNode := newRoot

	for level := uint32(1); level <= v.depth; level++ {
		copy(newNode.children[:], node.children[:])
		shifts := (v.depth - level) * bucketBits
		nodeIndex = (index >> shifts) & bucketMask
		nextNode := &vectorNode{}
		newNode.children[nodeIndex] = nextNode
		newNode = nextNode
		node = node.children[nodeIndex]
		if node == nil {
			node = &vectorNode{}
		}
	}

	copy(newNode.values[:], node.values[:])
	newNode.values[nodeIndex] = value

	return Vector{
		size:     v.size,
		capacity: v.capacity,
		depth:    v.depth,
		offset:   v.offset,
		root:     newRoot,
	}
}

func (v Vector) Get(index uint32) interface{} {
	if index >= v.size {
		panic("Out of bounds vector access")
	}

	index += v.offset

	node := v.root
	nodeIndex := index

	for level := uint32(1); level <= v.depth; level++ {
		shifts := (v.depth - level) * bucketBits
		nodeIndex = (index >> shifts) & bucketMask
		node = node.children[nodeIndex]
		if node == nil {
			return nil
		}
	}

	return node.values[nodeIndex]
}

func (v Vector) Append(value interface{}) Vector {
	appended := v.Resize(v.size + 1)
	return appended.Set(v.size, value)
}

func (v Vector) Resize(size uint32) Vector {
	offset := v.offset
	if size == 0 {
		offset = 0
	}

	capacity := v.capacity
	depth := v.depth
	root := v.root

	if capacity == 0 {
		capacity = bucketSize
		depth = 1
		root = &vectorNode{}
	}

	for size > capacity {
		capacity *= bucketSize
		depth++
	}

	return Vector{
		size:     size,
		capacity: capacity,
		depth:    depth,
		offset:   offset,
		root:     root,
	}
}

func (v Vector) Slice(start, end uint32) Vector {
	if end < start {
		panic("Invalid range")
	}
	if end == start || start >= v.size {
		return Vector{}
	}
	if end >= v.size {
		end = v.size
	}

	return Vector{
		size:     end - start,
		capacity: v.capacity,
		depth:    v.depth,
		offset:   start,
		root:     v.root,
	}
}

func (v Vector) Size() uint32 {
	return v.size
}
