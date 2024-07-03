---
title: Generic type `T any` cannot be `nil`
tagline: Golang code with unexpected behaviors
description: Unexpected golang behavior when trying to return a `nil` value for a generic type `T any`
---

# Generic type `T any` cannot be `nil`

The [generic type `T any`](https://go.dev/doc/go1.18#type-parameters) is a new feature in go1.18 that allows you to define a function that can accept any type, and return any type.

However, you can't return `nil` for a generic type `T any`, for example, this code would be invalid:

```go
func doSomething[T any]() T {
    return nil
}
```

# Error: cannot use nil as T value in return statement

The compiler will return an error like this: 

```text
cannot use nil as T value in return statement
```

At first glance, one might think it's unexpected, however notice that the generated code should be valid for specific types. The type `T` isn't really an `interface{}` type (`any`), but a placeholder for a specific type.

So if the compiler wanted to build the  `doSomething` function for a `string`, it would be like this:

```go
func doSomething() string {
    return nil
}
```

Which is now obviously invalid.

# How to return zero value for a generic type `T any`

If you want to return the zero value for a generic type `T any`, you can just declare a variable of type `T` and return it:

```go
package main

import "fmt"

func nothing[T any]() T {
	var zero T
	return zero
}

func main() {
	fmt.Printf("string: %q\n", nothing[string]())
	fmt.Printf("*int: %v\n", nothing[*int]())
}
```

If you [run this example in go.dev playground](https://go.dev/play/p/IxMnITpJXKT), it would print:

```text
string: ""
*int: <nil>
```


