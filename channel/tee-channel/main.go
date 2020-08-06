package main

func main() {
	//分割来自通道上的多个值，然后将它们发送到两个独立区域
	// 一个通道上接受一系列指令，同时记录操作日志

	// 这边还是需要一个orDone函数来做辅助的

	orDone := func(done <-chan struct{}, readStream <-chan interface{}) <-chan interface{} {
		ret := make(chan interface{})
		go func() {
			defer close(ret)
		loop:
			for {
				select {
				case <-done:
					break loop
				case val, ok := <-readStream:
					if ok == false {
						return
					}
					select {
					case <-done:
					case ret <- val:
					}
				}
			}
		}()
		return ret
	}

	tee := func(done <-chan struct{}, in <-chan interface{}) (<-chan interface{}, <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out2)
			defer close(out1)

			for val := range orDone(done, in) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					// 两个out送完了，就拿新的值来送
					// 重点在于置空，就可以让下一次这个堵塞
					select {
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					case <-done:
						// 其实不用return,因为一旦done了，外面这个range也结束了
						//return
					}
				}
			}
		}()
		return out1, out2
	}
}
