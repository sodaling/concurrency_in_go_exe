package main

import (
	"fmt"
	"sync"
)

func main() {
	// 启动两个，并且打印出两个正在睡觉，然后退出主go
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("go程1已结启动")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("go程2已结启动")
	}()
	wg.Wait()
	fmt.Println("两个go程都启动完毕了")
}
