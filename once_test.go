package webeh

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"fmt"
)

func TestOnce_Do(t *testing.T) {
	var value int64 = 0
	var o Once
	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Microsecond * (1 + time.Duration(rand.Int()%10)))
			o.Do(func() {
				atomic.AddInt64(&value, 1)
				time.Sleep(time.Microsecond * 200)
			})
			if value != 1 {
				t.Fatal(fmt.Sprintf("value %d is NOT equal to 1", value))
			}
		}()
	}
	wg.Wait()
	o.Reset()
	if value != 1 {
		t.Fatal(fmt.Sprintf("value %d is NOT equal to 1", value))
	}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.Do(func() {
				time.Sleep(time.Microsecond * 200)
				atomic.AddInt64(&value, -2)
			})
			if value != -1 {
				t.Fatal(fmt.Sprintf("value %d is NOT equal to -1", value))
			}
		}()
	}
	wg.Wait()
	if value != -1 {
		t.Fatal(fmt.Sprintf("value %d is NOT equal to -1", value))
	}
}

func TestOnce_Do2(t *testing.T) {
	var value int64 = 0
	var o1, o2 sync.Once
	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Microsecond * (1 + time.Duration(rand.Int()%10)))
			o1.Do(func() {
				atomic.AddInt64(&value, 1)
				time.Sleep(time.Microsecond * 200)
			})
			if value != 1 {
				t.Fatal(fmt.Sprintf("value %d is NOT equal to 1", value))
			}
		}()
	}
	wg.Wait()
	if value != 1 {
		t.Fatal(fmt.Sprintf("value %d is NOT equal to 1", value))
	}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o2.Do(func() {
				time.Sleep(time.Microsecond * 200)
				atomic.AddInt64(&value, -2)
			})
			if value != -1 {
				t.Fatal(fmt.Sprintf("value %d is NOT equal to -1", value))
			}
		}()
	}
	wg.Wait()
	if value != -1 {
		t.Fatal(fmt.Sprintf("value %d is NOT equal to -1", value))
	}
}

func BenchmarkOnce_Do(b *testing.B) {
	var value int64 = 0
	var o Once
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			o.Do(func() {
				atomic.AddInt64(&value, 1)
			})
		}
	})
}

func BenchmarkOnce_Do2(b *testing.B) {
	var value int64 = 0
	var o sync.Once
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			o.Do(func() {
				atomic.AddInt64(&value, 1)
			})
		}
	})
}
