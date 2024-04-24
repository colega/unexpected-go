---
title: `defer` does not wait for parallelized tests
tagline: Unexpected Golang's stuff
description: When calling `t.Parallel()` the inner tests may keep running even when the outer test function has finished.
---

# Defer does not wait for parallelized tests 

Consider the following test:

```go
func TestFoo(t *testing.T) {
	defer t.Logf("outer")

	t.Run("parallelized test", func(t *testing.T) {
		t.Parallel()
		t.Logf("inner")
	})
}
```

One could expect it to log `inner` first, and then `outer`. 

However, that's not the case, and [it would actually log them in the opposite order](https://go.dev/play/p/9LrrlLcSPR7).

# Why?

When running subtests using `t.Run()`, those run on separate goroutines. 
That allows them to be parallelized from the subtest itself by calling `t.Parallel()`.

When a sub-test calls `t.Parallel()`, the test runner pauses that goroutine and keep collecting the rest of the tests,
in order to track which ones are can be parallelilzed, and which one cannot be.

In order to account for all tests, the entire outer test function has to be executed, including the defer statement.

# Solution

The correct way of performing cleanup in the tests is to call `t.Cleanup()`, which will account for all the tests that are still pending to be finished.

The [following code would log `inner` first, and then `outer`](https://go.dev/play/p/0wr7Eqt7_Zb):

```go
func TestFoo(t *testing.T) {
	t.Cleanup(func() { t.Logf("outer") })

	t.Run("parallelized test", func(t *testing.T) {
		t.Parallel()
		t.Logf("inner")
	})
}
```

Even if you don't have parallelized tests today, don't use `defer` in your tests as they might cause subtle bugs in the future, resulting in long debugging hours.
