---
title: `os.Stdout` changes after invoking init() in examples
description: Stdout value may be unexpected in the examples
--- 

# `os.Stdout` changes after invoking init() in examples

Golang language has a cool feature: [testable
examples](https://blog.golang.org/examples), which basically runs your code as a
test, and asserts on the output you write into stdout.

Like this one:

```go
package main

import "fmt"

func ExampleSum() {
	fmt.Println(1 + 2)

	// Output:
	// 3
}
```

This is a test, it passes, it counts for the coverage, etc.

# Unexpected behavior

Consider this test:

```go
package main

import (
	"os"

	"log"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Don't log time, we can't assert on that
}

func ExampleLog() {
	log.Print("Log!")

	// Output:
	// Log!
}
```

We don't have a link to golang playground this time, but the result is:

```
$ go test ./...
Log!
--- FAIL: ExampleLog (0.00s)
got:

want:
Log!
FAIL
FAIL	github.com/cabify/product_go/testableexample	0.007s
```

So, we've printed _Log!_, but the test didn't get the _Log!_ it wanted.

Can we fix this? Changing the logger output in the example code:

```go
package main

import (
	"os"

	"log"
)

func ExampleLog() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0) // Don't log time, we can't assert on that

	log.Print("Log!")

	// Output:
	// Log!
}
```

Which passes:

```
$ go test ./...
ok  	github.com/cabify/product_go/testableexample	0.011s
```

# Why? Isn't init() supposed to be called when the package is initialized?

`init()` is actually called here, but golang's testimg framework changes the
value of os.Stdout after initializing the package:

```go
package main

import (
	"fmt"
	"os"
	"reflect"
)

var osStdout *os.File

func init() {
	osStdout = os.Stdout
	fmt.Println("Yes, we ran the init()")
}

func ExampleStdoutChanges() {
	fmt.Println(reflect.DeepEqual(osStdout, os.Stdout))
	// Output:
	// true
}
```

Prints:

```
$ go test ./...
Yes, we ran the init()
--- FAIL: ExampleStdoutChanges (0.00s)
got:
false
want:
true
FAIL
FAIL	github.com/cabify/product_go/testableexample	0.016s
```

This is probably done to intercept the `os.Stdout` for the example (notice that
the text we print from `init()` is not considered as `got:`). 
