package main

import (
	"sync"
	"sync/atomic"
)

func main() {
	//CountAtomic()Å
	CountMutex()
}

func CountAtomic() {
	var wg sync.WaitGroup
	count := int64(0)

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt64(&count, 1)
			wg.Done()
		}()
	}
	wg.Wait()
	//fmt.Println(count)
}
func CountMutex() {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	count := 0

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			defer mutex.Unlock()

			count++
			wg.Done()
		}()
	}
	wg.Wait()
	//fmt.Println(count)
}
