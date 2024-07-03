package rwmutex

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkRWMutex_FastRLock(b *testing.B) {
	for _, slowestReaderSpeed := range []time.Duration{10 * time.Millisecond, time.Second} {
		b.Run(fmt.Sprintf("slowest_reader=%s", slowestReaderSpeed), func(b *testing.B) {
			mtx := sync.RWMutex{}
			shared := 0

			done := make(chan struct{})
			wg := sync.WaitGroup{}
			for i := 0; i < runtime.GOMAXPROCS(0); i++ {
				// Fast readers, they don't lock for too long.
				wg.Add(1)
				go reader(&mtx, &wg, done, func() {
					if shared < 0 {
						panic("shared overflow")
					}
				})
			}

			wg.Add(1)
			go reader(&mtx, &wg, done, func() {
				time.Sleep(slowestReaderSpeed)
				if shared < 0 {
					panic("shared overflow")
				}
			})

			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-done:
						return
					case <-time.After(time.Millisecond):
						mtx.Lock()
						shared++
						mtx.Unlock()
					}
				}
			}()
			b.Cleanup(func() {
				close(done)
				wg.Wait()
			})

			b.RunParallel(func(pb *testing.PB) {
				sum := 0
				for pb.Next() {
					mtx.RLock()
					sum += shared
					mtx.RUnlock()
				}
				if sum < 0 {
					panic("sum overflow")
				}
			})
		})
	}
}

func reader(mtx *sync.RWMutex, wg *sync.WaitGroup, done chan struct{}, f func()) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			mtx.RLock()
			f()
			mtx.RUnlock()
		}
	}
}
