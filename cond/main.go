package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 简单来说是等通知，一个生产者ok了，通知多个消费者来拿
	// 他和一般的channel堵塞不一样的是，如果不用close channel的方法来通知的话，这个通知只能通知到一个go程。而如果只用close来通知的话
	// 如果我需要消费者轮流消费，而不是一下子都醒了，就会得多用一个锁。很麻烦
	// 就比如我生产者一下子生产了5个苹果，但是有10个消费者在等着。就用得上这个了。
	// 现在这里有个案例，假设我们有一个固定长度为2的队列，并且我们要将10个元素放入队列中，我们希望一有空间就能放入，所以在队列中有空间时需要立刻通知。
	// 临界资源还是要用锁锁住的
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("queue已经remove了一个")
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			// 这边加上判断的原因是，因为是一个个唤醒的，所有有可能后面唤醒的时候已经不符合唤醒的条件了，所以如果看到条件不对就得重新继续wait
			c.Wait()
		}
		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
	}
	fmt.Printf("还有%d个", len(queue))
}
