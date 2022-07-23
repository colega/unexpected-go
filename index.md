---
title: Unexpected Go
tagline: Golang code with unexpected behaviors
description: Showcase of unexpected golang behaviors
---

This small site is oriented to developers learning [golang](https://golang.org)
that want to see the _unexpected_, at least for a newcomer, behaviors in the
language.

This examples expect that reader is familiar with the language, if you're not
familiar with it, you should complete the [Go by
Example](https://gobyexample.com/) tutorial.

# Examples 

The WIP list of examples is (with no specific order):
 - [`nil` errors that are non-nil `error`s](nil-errors-that-are-non-nil-errors.html)
 - [`time.RFC3339` constant is not a valid `RFC3339` timestamp](time-rfc3339-is-not-a-valid-rfc3339.html)
 - [`os.Stdout` changes after `init()` in testable examples](os-stdout-changes-after-init-in-examples.html)
 - [`Write([]byte) (int, error)` should not retain the slice](write-should-not-retain-the-slice.html)
 - [`if err := recover(); err != nil` does not recover from all panics](recovering-from-nil-panics.html)
 - [There's no `Min(int, int)` int function](theres-no-min-function.html)
 - [Goroutines on `for` loop values](goroutines-on-loops.html)
 - [`json.Unmarhsal` mreges structs and slices, but not map values](json-unmarshal-merges-structs-but-not-map-values.html)

# Grab some popcorn
![Unexpected gopher](gopher.png)
   
# Disclaimer 
   
This examples are not trying to prove that golang is _wrong_, it's just
sharing stuff that someone who has just started learning the language doesn't
expect.

