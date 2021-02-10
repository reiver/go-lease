# go-lease

Package **lease** provides a safer way of working with mutexes in the Go programming language.

It is an alternative to the Go programming language's built-in `sync.Mutex` type in the `sync` package.

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-lease

[![GoDoc](https://godoc.org/github.com/reiver/go-lease?status.svg)](https://godoc.org/github.com/reiver/go-lease)

## Basic Usage

Basic usage is as follows.
```go
var c lease.Type

err := c.Lease(func(){
	//@TODO
})

switch casted := err.(type) {
case lease.Timedout:
	//@TODO
default:
	//@TODO
}
```
