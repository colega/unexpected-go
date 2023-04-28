---
title: Nil errors that are non-nil errors
description: Difference between nil interface value and interface holding a nil value
---

# Nil errors that are non-nil errors

Consider [the following code](https://play.golang.org/p/j_W4udUC72o):

```go
package main

import "fmt"

type E struct {}

func (*E) Error() string { return "error" }

func f() error {
    var err *E = nil
    return err
}

func main() {
    if err := f(); err != nil {
        panic(err)
    }
    fmt.Println("success")
}
```

You can easily check that although `f() error` always returns `nil`, `f() != nil` so this code panics.

# Why?

[An interface in Go is a tuple of (type, value)](https://golangbot.com/interfaces-part-1/).
In the following code:

```go
var err error
```

`err` has `nil` value since we didn't assign anything to it, but in:

```go
var (
    specific *E
    err error = specific
)
```

We're setting the err's value to a `(*E, nil)` tuple, there's a type, but it doesn't point to any value.

You can check that out [in the following example](https://play.golang.org/p/IK1AWaO4HRX):

```go
package main

import "fmt"

type Printer interface{ Print() }

type StringPrinter string

func (s *StringPrinter) Print() {
	if s == nil {
		fmt.Println("nil value")
		return
	}
	fmt.Println(*s)
}

func main() {
	var (
		sp      *StringPrinter
		printer Printer = sp
	)

	printer.Print()
}
```

Since printer already has a type, we can already call pointer receiver functions on it, although it's value is nil.
Obviously you [can't do that](https://play.golang.org/p/A-A1aDZejR9) if you don't specify the type, because there's nowhere to call.
