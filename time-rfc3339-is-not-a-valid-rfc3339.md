---
title: time.RFC3339 is not a valid time.RFC3339 format
description: Just a stupid curiosity that will never affect your productivity
----

# Introduction

You may know that there's a `time.Parse(layout, value string) (time.Time, error)` function and that there are some pre-defined layouts in the `time` package like `time.RFC3339`, for instance.

So, so on one sunny day you decide to use some default value for a [RFC3339](https://tools.ietf.org/html/rfc3339) timestamp in one of your tests, and first value that comes to your mind is the `time.RFC3339` constant which value is 

```go
const(
    // ...
    RFC3339     = "2006-01-02T15:04:05Z07:00"
    // ...
)
```

So if you don't pay a lot of attention, you may think that you can parse `time.RFC3339` with `time.RFC3339` as the layout, [but you actually can't](https://play.golang.org/p/0MRMopw14gG):

```go
package main

import (
	"time"
)

func main() {
	_, err := time.Parse(time.RFC3339, time.RFC3339)
	if err != nil {
		panic(err)
	}
}
```

Outputs:
```
panic: parsing time "2006-01-02T15:04:05Z07:00": extra text: 07:00
```

# Why? 

Because the layout itself is not a valid timestamp, since it includes both UTC (`Z`) and timezone definition (`07:00`), while the timestamp should include only one of them.

