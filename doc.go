/*
Package immutable provides immutable containers for building
applications that handle concurrency without locks.

The implementation tries to be as straightforward and minimal as possible,
while still being performant enough for everyday use.

Basics

The containers in this package return a new version of themselves
as the result of any updating operation. This makes the data structure as a
whole immutable. There is no way to change existing instances.

Iterating during modification for example, is not an issue.
It simply cannot happen.

The containers hold values of type interface{}.

Copying

Making copies of an immutable container is cheap and independent of the
current size.

Copying is not atomic, however, just like as there is no guarantee for atomic
assignment of pointers.

Use channels as usual to communicate, or atomic.Value if you really need to.

Performance

Comparing these containers to built-in constructs in Go or other mutable
containers is a bit of an apples to oranges comparison.

Modifying a container creates a new instance, which requires memory allocations.
There is no need for mutexes to make the containers safe for use in multiple
concurrent threads, however.

And the Vector for example, is sparse, so growing from 1 to a million addressable
elements is much faster than the same operation for slices.

Check out the benchmarks for some example comparisons of performance.

Basic use

A zero-valued container is empty and ready to use.

	var myMap Map
	myMap = myMap.Set("thing", 42)

	a := Vector{}.Resize(4711)
	b := a.Set(1000, "hallå!")
*/
package immutable
