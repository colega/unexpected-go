---
title: `clear()` removes map keys but keeps slice indexes
description: The built-in `clear()` function removes all the keys from a map, but it does not remove the indexes from a slice.  
---

#  `clear()` built-in function removes map keys but keeps slice indexes

[The documentation](https://pkg.go.dev/builtin#clear) for the `clear()` built-in function that was introduced in go1.20.0 states this clearly:

> For maps, clear deletes all entries, resulting in an empty map. For slices, clear sets all elements up to the length of the slice to the zero value of the respective element type.

While documented and correct, this still falls into "unexpected" behaviors category.

For example [this code](https://play.golang.com/p/ui4ayQmU8b-):

```go
package main

import (
	"fmt"
)

func main() {
	m := map[string]int{"one": 1, "two": 2}
	s := []int{1, 2}

	fmt.Println("Original:")
	fmt.Printf("m=%#v (len %d)\n", m, len(m))
	fmt.Printf("s=%#v (len %d)\n", s, len(s))
	fmt.Println()

	clear(m)
	clear(s)
	fmt.Println("Cleared:")
	fmt.Printf("m=%#v (len %d)\n", m, len(m))
	fmt.Printf("s=%#v (len %d)\n", s, len(s))
	fmt.Println()
}
```

Prints:

```text
Original:
m=map[string]int{"one":1, "two":2} (len 2)
s=[]int{1, 2} (len 2)

After clear():
m=map[string]int{} (len 0)
s=[]int{0, 0} (len 2)
```

As you can see, the length of the map is zeroed, while the length of the slice remains the same, but the values are zeroed.

# Why?

Consider reading [the spec issue](https://github.com/golang/go/issues/56351) on the official golang repository for more context.

One of the reasons might be that you [can't always remove all keys from a map using the `delete()` function](you-cant-remove-a-nan-key-from-a-map-without-clearing-it.html), while there's a builtin syntax for changing a slice length.

# Extra

It's also worth mentioning that `clear()` clears up to the length of the provided slice, not up to it's capacity, [for example](https://play.golang.com/p/4p8gkV01x-W):

```go
package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2}
	fmt.Println("Original:")
	s = s[:1]
	clear(s)
	fmt.Println("Cleared:")
	fmt.Println(s)
	s = s[:2]
	fmt.Println("Extended:")
	fmt.Println(s)
}
```

Outputs:

```text
Original:
Cleared:
[0]
Extended:
[0 2]
```