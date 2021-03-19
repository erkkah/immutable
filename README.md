# Immutable datastructures for Golang

This package provides immutable Map and Vector implementations for Golang.
These data structures are immutable in the sense that all operations that
modify the structure return a modified copy of the structure.

Copying each datastructure is fast, but there is no guarantee from the
runtime that it will be atomic. Use atomic.Value or similar.

Once you have a copy, it is your immutable copy, and will never change.

[API docs](https://pkg.go.dev/github.com/erkkah/immutable)
