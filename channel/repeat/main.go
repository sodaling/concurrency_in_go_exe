package main

import (
	"fmt"
	"time"
)

func main() {
	// 一个简单的repeat生成器
	repeat := func(done <-chan struct{}, values ...int) <-chan int {
		ret := make(chan int)
		go func() {
			// 一定要记得close
			defer close(ret)
			for {
				for _, i := range values {
					select {
					case ret <- i:
					case <-done:
						return
					}
				}
			}
		}()
		return ret
	}
	temp := []int{1, 2, 3, 4, 5, 6}
	doneTwo := make(chan struct{})
	go func() {
		fmt.Println("开始睡觉")
		time.Sleep(1 * time.Second)
		fmt.Println("醒来了")
		close(doneTwo)
	}()
	for i := range repeat(doneTwo, temp...) {
		fmt.Printf("%d输出啦\n", i)
	}
	fmt.Println("没了")

	take := func(done <-chan struct{}, valueStream <-chan int, num int) <-chan int {
		ret := make(chan int)
		go func() {
			defer close(ret)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case ret <- <-valueStream:
					// 注意要两个<-，一个取出一个灌入
				}
			}
		}()
		return ret
	}
	doneThree := make(chan struct{})
	go func() {
		fmt.Println("开始睡觉,再次")
		time.Sleep(1 * time.Second)
		fmt.Println("醒来了,再次")
		close(doneTwo)
	}()
	for i := range take(doneThree, repeat(doneThree, temp...), 10) {
		fmt.Printf("%d输出啦(再次)\n", i)
	}
	fmt.Println("真没了")
}
