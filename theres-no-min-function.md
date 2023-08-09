---
title: There is no math.Min(int, int) int function
tagline: Unexpected Golang's stuff
description: There is no min(int, int) int, or max function in math, and there was no builtin function before go1.21
---

# TL;DR: There's no `math.Min(int, int) int` function in golang, but there's a builtin since go1.21

Just as easy as that, if you check the [math package](https://golang.org/pkg/math/) you won't find a `func Min(a, b int) int`, neither a `Max` for ints, or for int64, etc.

However, [since go1.21](https://go.dev/doc/go1.21) there are [min and max builtin functions](https://tip.golang.org/ref/spec#Min_and_max).

# What happened before?

[Go does not support overloading of methods and operators](https://golang.org/doc/faq#overloading). 
In many other languages you can define two functions as: `min(a, b int)` and `min(a, b float64)` but if you try this in go your program won't compile. 
Since defining a bunch of different functions, one for each numerical type (`MinInt`, `MinInt64`, `MinFloat64`, ...), 
is kind of messy and go is aimed at simplicity there is only `Min(a, b float64)`.

# Solutions

## Casting ints into floats

There's one for floats, and you'll have to deal with that, either [casting your ints to floats](https://play.golang.org/p/78mrx3gqttq):

```go 
package main

import (
	"fmt"
	"math"
)

func main() {
	a := 1
	b := 2
	m := math.Min(float64(a), float64(b))
	
	fmt.Println(m)
}
```

(Note: casting ints into floats may result in loss of precision)

## Custom specifics function 

Defining [your own min function everywhere](https://play.golang.org/p/JppKTRa34pF) is an option:

```go
package main

import (
	"fmt"
)

func main() {
	a := 1
	b := 2
	m := min(a, b)

	fmt.Println(m)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
```

## Generics since Go 1.18

[Go 1.18 introduced support for _generics_](https://tip.golang.org/doc/go1.18), and now [a single function can be defined](https://go.dev/play/p/R4DNAgKyWo2) for all the _comparable types_:

```go
package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println(min(1, 2))
	fmt.Println(min(1.5, 2.7))
}
```

You can define it for yourself, or use one of the multiple generic packages out there but [remember the proverb](https://go-proverbs.github.io/): [a little copying is better than a little dependency](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s).
