---
title: `io.Writer`'s `Write` should not retain the slice
tagline: Golang code with unexpected behaviors
description: Unexpected golang behavior when a Writer stores the reference to the provided bytes.
---

# `io.Writer`'s `Write` should not retain the slice

One should always pay special attention about the lifecycle of the data structures when they're shared between different entities, i.e., they're pointers.
Usually when the referenced value can change, some explicit concurrency mechanisms are offered, or the lifecycle is clearly documented.
In the case of the `Writer` interface provided by the standard `io` package, the latter is applied, and [the documentation][doc] clearly states that _implementations must not retain p._, however due to the simplicity of the function one may overlook the documentation and suffer from that.

This section explains how retaining the `[]byte` slice provided to `Write([]byte) (int, error)` can bite you.

# Code

Imagine yourself writing a dummy buffered writer implementation, because you don't want your program to be slowed down by your logging system, one of the simpliest approaches would [look like this][broken code]:

```go
package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type bufWriter struct {
	output io.Writer
	buffer chan []byte
}

func (bw bufWriter) Start() {
	go func() {
		for b := range bw.buffer {
			bw.output.Write(b)
		}
	}()
}

func (bw bufWriter) Write(b []byte) (int, error) {
	bw.buffer <- b
	return len(b), nil
}

func main() {
	writer := bufWriter{os.Stdout, make(chan []byte, 10)}
	writer.Start()

	for i := 0; i < 10; i++ {
		fmt.Fprintf(writer, "%d\n", i)
	}

	time.Sleep(time.Second)
}
```

One would expect this code to print numbers from 0 to 9, however, the output will be most likely a bunch of nines.

# Why? 

The cause of this, apart obviously of our violation of the contract of `Write` method are some optimizations performed by `fmt` package.
`fmt.Fprintf` calls `fmt.newPrinter()` function which doesn't actually always create a new `pp`, but uses a `sync.Pool` of them instead.

```go
// newPrinter allocates a new pp struct or grabs a cached one.
func newPrinter() *pp {
	p := ppFree.Get().(*pp)
	p.panicking = false
	p.erroring = false
	p.wrapErrs = false
	p.fmt.init(&p.buf)
	return p
}
```

It uses a `buffer` to build the formatted string, then passes that buffer to our `Write()` implementation and then calls `pp.free()` which resets some values on the printer and puts it back to the pool so it can be reused for the next formatting.

When we call `fmt.Fprintf` again, same printer is being reused, and the buffer, although it's `[]byte` slice is not 0, still points to the same underlying array, so we end up writing on the same data.
And here's where the issue comes: we still didn't read that slice, and it's being overwritten already, so when we read it, it's already the latest value that have been formatted.

# Solution

Well, as we already mentioned we should just follow [the documentation][doc] and avoid retaining the `[]byte` slice we've been given. Instead of that, we should copy it's contents before returning from the `Write()` function, [like this][fixed code]:

```go
func (bw bufWriter) Write(b []byte) (int, error) {
	cp := make([]byte, len(b))
	for i := range b {
		cp[i] = b[i]
	}
	bw.buffer <- cp
	return len(cp), nil
}
```

## Optimization

Since allocating a new slice each time may be expensive (remember, never assume, measure first), you may want to use a [`sync.Pool` combined with a `bytes.Buffer`][optimized code].

[doc]: https://golang.org/pkg/io/#Writer 
[broken code]: https://play.golang.org/p/8XFwXMRh-oe
[fixed code]: https://play.golang.org/p/c9eOepzt8G3
[optimized code]: https://play.golang.org/p/yeJAhwLsYZQ
