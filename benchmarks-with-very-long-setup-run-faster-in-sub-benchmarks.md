---
title: Benchmarks with very long setup run faster in sub-benchmarks
description: Benchmark running time can be optimized by wrapping them into trivial sub-benchmarks.
--- 

# Introduction

Golang language has [a cool feature that allows you writing Benchmarks](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) with just the standard toolkit.

For example, you can easily benchmark how long does the `time.Now()` call take by writing:

```go
func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now() // We're benchmarking this
	}
}
```

And running:

```
$ go test -run=X -bench=. -count=5 -cpu=1 .
goos: linux
goarch: amd64
pkg: github.com/grafana/mimir
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
BenchmarkTest 	46728801	        25.79 ns/op
BenchmarkTest 	46436253	        25.55 ns/op
BenchmarkTest 	46213618	        25.67 ns/op
BenchmarkTest 	46396131	        25.94 ns/op
BenchmarkTest 	46949673	        25.54 ns/op
PASS
ok  	github.com/grafana/mimir	7.087s
```

# Why benchmarks with a slow setup run so slow

However, sometimes you need to setup a heavy environment to run the benchmarks (fill a database for example).
In that case you will need to call `b.ResetTimer()` to avoid accounting for that setup in your benchmarked call:

```go
func BenchmarkTest(b *testing.B) {
	// Start of a very long setup
	time.Sleep(10 * time.Second)
	// End of a very long setup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now() // We're benchmarking this and it needs the setup
	}
}
```

However, if you run that benchmark you will quickly realize that even though the results are correct (notice the samea 25 ns/op result), the time it took to run the benchmark is quite big:

```
$ go test -run=X -bench=. -count=5 -cpu=1 .
goos: linux
goarch: amd64
pkg: github.com/grafana/mimir
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
BenchmarkTest 	45376972	        25.68 ns/op
BenchmarkTest 	39959161	        25.96 ns/op
BenchmarkTest 	46036717	        26.07 ns/op
BenchmarkTest 	44628370	        26.12 ns/op
BenchmarkTest 	46740477	        25.44 ns/op
PASS
ok  	github.com/grafana/mimir	266.654s
```

This happens because golang doesn't really know how the benchmarked operation will take, so it tries to estimate by running the benchmark with smaller `b.N` values first.
This really means that the entire setup is ran multiple times, you can check that by adding some logs to the benchmark:

```go
func BenchmarkTest(b *testing.B) {
	b.Log("Start of a very long setup")
	time.Sleep(10 * time.Second)
	b.Log("End of a very long setup")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		time.Now() // We're benchmarking this, and it needs the setup
	}
}
```

And see how the input logs that multiple times before even reaching the results of the first run:

```
go test -run=X -bench=. -count=5 -cpu=1 .
goos: linux
goarch: amd64
pkg: github.com/grafana/mimir
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
BenchmarkTest 	46137392	        25.63 ns/op
--- BENCH: BenchmarkTest
    bench_test.go:9: Performing a very long setup
    bench_test.go:11: End of a very long setup
    bench_test.go:9: Performing a very long setup
    bench_test.go:11: End of a very long setup
    bench_test.go:9: Performing a very long setup
    bench_test.go:11: End of a very long setup
    bench_test.go:9: Performing a very long setup
    bench_test.go:11: End of a very long setup
    bench_test.go:9: Performing a very long setup
    bench_test.go:11: End of a very long setup
	... [output truncated]

... rest of the output omitted
```

# How to make benchmarks with a slow setup run faster

There's a trick to make the setup run only once, which consists of wrapping the benchmarked code in a sub-benchmark `b.Run(...)`:

```go
func BenchmarkTest(b *testing.B) {
	b.Log("Start of a very long setup")
	time.Sleep(10 * time.Second)
	b.Log("End of a very long setup")

	b.Run("benchmark", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			time.Now() // We're benchmarking this, and it needs the setup
		}
	})
}
```

This allows Go's testing framework to re-run just the benchmark itself, running the setup only once:

```
go test -v -run=X -bench=. -count=5 -cpu=1 .
goos: linux
goarch: amd64
pkg: github.com/grafana/mimir
cpu: 11th Gen Intel(R) Core(TM) i7-11700K @ 3.60GHz
BenchmarkTest
    bench_test.go:9: Start of a very long setup
    bench_test.go:11: End of a very long setup
BenchmarkTest/benchmark
BenchmarkTest/benchmark         	46163953	        25.70 ns/op
BenchmarkTest/benchmark         	45692376	        25.81 ns/op
BenchmarkTest/benchmark         	46491103	        25.64 ns/op
BenchmarkTest/benchmark         	44900762	        25.61 ns/op
BenchmarkTest/benchmark         	47827808	        25.54 ns/op
PASS
ok  	github.com/grafana/mimir	16.068s
```

Hopefully will save you some precious time.
