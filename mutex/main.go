package main

import (
	"fmt"
	"sync"
)

func main() {
	// 设定一个数，5个go程序去加，5个go程去减
	var count int
	var lock sync.Mutex
	countIncrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("the count is %d now\n", count)
	}

	countDecrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("the count is %d now\n", count)
	}

	var wg sync.WaitGroup
	// 这边加
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			countIncrement()
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
			countDecrement()
		}()
	}
	wg.Wait()
	fmt.Printf("the count now is %d\n", count)
}
