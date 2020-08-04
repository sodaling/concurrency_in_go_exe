package main

import "fmt"

func main() {
	// 值得注意的是，我们在chanOwer的此法范围内实例化chan，这和将导致通道的写入操作范围被限制在
	// 它的下面定义的闭包中。换句话说，它限制了这个通道的写入使用范围，以防止其他goroutine写入它
	chanOwer := func() <-chan int {
		ret := make(chan int, 5)
		go func() {
			defer close(ret)
			for i := 0; i < 5; i++ {
				ret <- i
			}
		}()
		return ret
	}
	// 在这里，我们接受到一个只读的通道，我们传递给消费者，消费者只能从中读取信息
	// 同时这里收到一个int通道的只读副本，通过声明该函数的唯一用法是读取访问，我们将通道用法限制为只读
	consumer := func(results <-chan int) {
		for i := range results {
			fmt.Println(i)
		}
	}
	consumer(chanOwer())
}
