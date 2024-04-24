---
title: You can't remove a `NaN` key from a golang map without clearing it
description: Since `NaN` is not equal to itself, you can't remove it from a map by calling the `delete()` function.  
---

#  You can't remove a `NaN` key from a map without clearing it

Consider the following [example code](https://play.golang.com/p/-KcUKurzvJW):

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	m := map[float64]string{
		1:          "one",
		math.NaN(): "not a number",
	}
	
	fmt.Println("Original: ", m)
	delete(m, 1)
	delete(m, math.NaN())
	fmt.Println("After deleting: ", m)
}

```

It will output:

```text
Original:  map[NaN:not a number 1:one]
After deleting:  map[NaN:not a number]
```

As you can see, the `NaN` value is not removed from the map.

# Why?

`NaN` values [are not equal to themselves](https://play.golang.com/p/LwPyHBj2eFp), so the `delete()` function can't find the key to remove it.

# Solution

You can completely clear the map by calling the [`clear()` built-in function](https://pkg.go.dev/builtin#clear) that was introduced in go1.21.0, or alternatively you can build a new map by iterating the original one, and skipping the `NaN` keys by using the `math.IsNaN()` function:

# Extra

A corollary of `NaN` keys not being equal to other `NaN` keys is that you can have [as many `NaN` keys as you want](https://play.golang.com/p/vK2qM2RnJd0):

```go
package main

import (
	"fmt"
	"math"
)

func main() {
	m := map[float64]string{
		math.NaN(): "first",
		math.NaN(): "second",
	}

	fmt.Println("Original: ", m)
	m[math.NaN()] = "third"
	fmt.Println("Modified: ", m)
}
```

Prints:

```text
Original:  map[NaN:second NaN:first]
Modified:  map[NaN:third NaN:second NaN:first]
```