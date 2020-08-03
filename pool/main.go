package main

import (
	"fmt"
	"sync"
)

func main() {
	//myPool := &sync.Pool{New: func() interface{} {
	//	fmt.Println("Creating new instance.")
	//	return struct {
	//	}{}
	//},
	//}
	//myPool.Get()
	//instance := myPool.Get()
	//myPool.Put(instance)
	//// 只会输出两个新建语句
	//myPool.Get()

	// 实例化一个元素时候，给它一新的元素，这个元素应该是线程安全的
	// 当你从Get获得一个实例的时候，不要假设你接受到的对象状态
	// 当你从池咋红获得实例的时候，一定要记得put回去。通常用defer
	// 池内的元素必须大致上均匀的

	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			// 返回切片指针地址
			return &mem
		},
	}
	for i := 0; i < 4; i++ {
		// 先扩展到4kb
		calcPool.Put(calcPool.New())
	}
	// 和结尾对比看看，到底多创建了多少次
	fmt.Printf("%d calcultors were created.", numCalcsCreated)
	numWorker := 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}
	wg.Wait()
	//可以看到最多也就创建12次
	fmt.Printf("%d calcultors were created.", numCalcsCreated)
}
