---
title: There's no min(int, int) int function
tagline: Unexpected Golang's stuff
description: There's no min(int, int) int, or max function in math
---

# TL;DR

Just as easy as that, if you check the [math package](https://golang.org/pkg/math/) you won't find a `func Min(a, b int) int`, neither a `Max` for ints, or for int64, etc.

# Why?

[Go does not support overloading of methods and operators](https://golang.org/doc/faq#overloading). In many other languages you can define two functions as: `min(a, b int)` and `min(a, b float64)` but if you try this in go your program won't compile. Since defining a bunch of different functions, one for each numerical type (`MinInt`, `MinInt64`, `MinFloat64`, ...), is kind of messy and go is aimed at simplicity there is only `Min(a, b float64)`.

# Solutions

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

or defining [your own min function everywhere](https://play.golang.org/p/JppKTRa34pF):

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
