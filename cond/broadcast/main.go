package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func (but *Button) Subscribe(fun func()) {
	var tempWg sync.WaitGroup
	tempWg.Add(1)
	go func() {
		tempWg.Done()
		// 注意先要加锁，因为wait会解锁
		but.Clicked.L.Lock()
		defer but.Clicked.L.Unlock()
		but.Clicked.Wait()
		fun()
	}()
	tempWg.Wait()
	fmt.Println("订阅成功")
}

func main() {
	button := Button{sync.NewCond(&sync.Mutex{})}
	var wg sync.WaitGroup
	wg.Add(3)
	printOne := func() {
		fmt.Println("printOne")
		wg.Done()
	}
	printTwo := func() {
		fmt.Println("printTwo")
		wg.Done()
	}
	printThree := func() {
		fmt.Println("printThree")
		wg.Done()
	}
	button.Subscribe(printOne)
	button.Subscribe(printTwo)
	button.Subscribe(printThree)

	// 全部唤醒
	button.Clicked.Broadcast()

	wg.Wait()
}
