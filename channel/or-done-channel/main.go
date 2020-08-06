package main

import "fmt"

func main() {
	// 遍历一个channel时候，有时候需要随时退出。需要用select来分钟我们的读取操作和done通道
	orDone := func(done <-chan struct{}, c <-chan interface{}) <-chan interface{} {
		ret := make(chan interface{})
		go func() {
			defer close(ret)
		loop:
			for {
				select {
				case <-done:
					// 其实我觉得直接return也可以
					break loop
				case maybeVal, ok := <-c:
					if ok == false {
						return
					} else {
						select {
						case <-done:
						// 这边不许要做什么，跑完就可以了
						case ret <- maybeVal:
						}
					}

				}
			}
		}()
		return ret
	}
	// 这个就不掩饰了，反正就是可以一个range搞完遍历和中止
	fmt.Println(orDone)
}
