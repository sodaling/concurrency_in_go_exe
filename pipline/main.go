package main

import "fmt"

func main() {
	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		// 把一个整数切片转化成int类型的chan
		ret := make(chan int)
		go func() {
			defer close(ret)
			for _, i := range integers {
				select {
				case <-done:
					return
				case ret <- i:
				}
			}
		}()
		return ret
	}
	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		ret := make(chan int)
		go func() {
			defer close(ret)
			// 如果简单用range的话，会不会卡死。得保证之前的不卡死
			// 但是前面的流保证了，一旦前面的通过done关了，intStream也会跟着关闭，所以应该不会卡死
			for i := range intStream {
				select {
				case <-done:
					return
				case ret <- i * multiplier:
				}
			}
		}()
		return ret
	}
	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		ret := make(chan int)
		go func() {
			defer close(ret)
			for i := range intStream {
				select {
				case <-done:
					return
				case ret <- i + additive:
				}
			}
		}()
		return ret
	}
	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	for i := range add(done, multiply(done, intStream, 2), 2) {
		// 重点在于管道的每个阶段都在同时执行，任何阶段只需要等待其输出
		// 简单来说就是先*2然后+2
		// done可以保证程序干净退出而且不泄露go程
		// 允许每个阶段相互独立执行一段时间

		//关闭通道时候会怎么影响管道的:
		//1.对传入通道进行遍历的时候，当输入通道关闭时候，遍历操作将退出
		//2.发送操作与done通道共享select语句
		//保证不会堵死
		fmt.Println(i)
	}
}
