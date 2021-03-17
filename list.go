package immutable

// List is an immutable list with copy-on-write semantics.
// Modifying the list returns a new list instance.
// Since the list is immutable, it is safe to use from
// multiple concurrent threads without locks or other
// synchronization.
//
// The zero List is empty and ready for use.
type List struct {
	size uint32
}

func (l List) Set(index int, value interface{}) List {

}

func (l List) Append(value interface{}) List {

}

func (l List) Resize(int size) List {

}

func (l List) Slice(int start, int end) List {

}

func (l List) Get(index int) interface{} {

}

func (l List) Size() int {

}
