---
title: Recovering from nil panics 
description: A panic value can be nil
---

# Introduction

After finishing the Golang tutorial, one expects that this code will recover from any panic:

```go
func main() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Println("Recovered from panic", err)
        }
    }()

    // ...
}
```

However it will not, [if the panic value is `nil`](https://play.golang.org/p/MJxoRVNak7V):

```go
func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Recovered from panic", err)
		}
	}()

	panic(nil)
}
```

# Workarounds

[_bradfitz_](https://github.com/bradfitz) mentions as [the proper way](https://github.com/golang/go/issues/25448) something [like this](https://play.golang.org/p/1uKKlQg3t1A):

```go
func main() {
	panicked := true
	defer func() {
		if err := recover(); err != nil || panicked {
			fmt.Println("Recovered from panic", err)
		}
	}()

	panic(nil)
	panicked = false
}
```

An alternative approach would also be to wrap your calls with some kind of function that would transform nil panics into non-nil panics:

```go
func runWithoutNilPanics(f func()) {
    panicked := true
	defer func() {
		if err := recover(); err != nil {
		    panic(err)
		} else if panicked {
		    panic(fmt.Errorf("panicked with nil panic"))
		}
	}()

	f()
}
``` 
