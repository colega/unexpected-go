---
title: "Goroutines in loops  Unexpected Go"
tagline: Golang code with unexpected behaviors
description: Unexpected golang behavior when running goroutines with values from a loop
---

# Introdution

This is actually documented in [CommonMistakes](https://github.com/golang/go/wiki/CommonMistakes) 
wiki page of the official golang's GitHub repository (notice, this is the only common mistake documented
there as this page is being written) but since its visibility is not good enough, it's always a
good start. 

# Code

Consider the following [code](https://play.golang.org/p/5dyYBH6f05X):

```go
package main

import "fmt"

func main() {
	ch := make(chan int)
	slice := []int{0, 1}

	for _, v := range slice {
		go func() { ch <- v }()
	}

	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

```

One may expect this to print:

```
0
1
```

But actually this code is not deterministic, and in most cases will just print
```
1
1
````

# Why?

When iterating using `i, v := range whatever`, the `i` and `v` variables 
are defined only once for the scope, and then their values are overwritten.
Since the goroutines reference the same variable all the time, but they are 
_probably_ executed after the `for` loop is finished, they all access the same 
last value of it: `1`.

# Related

## `for ;;`
This can be also unexpected [if we do](https://play.golang.org/p/3EcWDaPKC8h):

```go
package main

import "fmt"

func main() {
	ch := make(chan int)

	for i := 0; i < 2; i++ {
		go func() { ch <- i }()
	}

	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

Where the output is:
```
2
2
```

Although our goroutines should have never been executed with `i >= 2` since `i < 2` is the stop condition for the loop.

## Receiver call

As a corolary, but still worth mentioning, this also applies to receiver functions, as a [simple example](https://play.golang.org/p/3JxDLDF75sO):

```go
package main

import (
	"fmt"
	"sync"
)

type S string

func (s S) Print() {
	fmt.Println(s)
}

func main() {
	slice := []S{"a", "b"}

	wg := sync.WaitGroup{}
	wg.Add(2)

	for _, s := range slice {
		go func() {
			s.Print()
			wg.Done()
		}()
	}

	wg.Wait()
}
```

Prints:

```
b
b
```
