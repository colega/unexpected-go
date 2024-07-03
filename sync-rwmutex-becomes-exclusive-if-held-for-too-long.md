---
title: sync.RWMutex becomes exclusive if held for too long
tagline: Golang code with unexpected behaviors
description: No RLock() can suceeed if there's a Lock() waiting, and the Lock() will wait for all RLock() to finish
---

# `sync.RWMutex` becomes an exclusive if held for too long

One of the most common synchronization primitives in golang is the [`sync.RWMutex`](https://pkg.go.dev/sync#RWMutex) which is a reader-writer mutual exclusion lock.

The `sync.RWMutex` has two main methods: `RLock()` and `Lock()`, the first one is used to allow multiple readers to access the resource, while the second one is used to allow only one writer to access the resource.

However, if a `RLock()` is held for too long, it will block any `Lock()` from being executed, and the waiting `Lock()` will wait for all `RLock()` to finish.

So, as a corollary, if you have `Lock()` calls happening more frequently than the longest `RLock()` is being held, the mutex will behave as an exclusive lock and in the best case you'll have bottleneck in your code.

