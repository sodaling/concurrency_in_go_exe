package main

import (
	"fmt"
	"time"
)

func main() {
	// 合并一个或者多个done通道到一个done通道中，该通道在其中任何一个组件通道关闭的时候关闭
	// 简单来说就是利用了递归解决这个问题
	// 要声明了才能调用
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		ret := make(chan interface{})
		go func() {
			defer close(ret)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], ret)...):
					// 上面要拆包，从第三个开始继续往外新建go程

				}
			}
		}()
		return ret
	}
	sig := func(after time.Duration) <-chan interface{} {
		// sleep一段时间后关闭返回
		ret := make(chan interface{})
		go func() {
			defer close(ret)
			time.Sleep(after)
		}()
		return ret
	}
	start := time.Now()
	<-or(sig(1*time.Second), sig(1*time.Hour), sig(1*time.Minute))
	fmt.Printf("%v秒后结束", time.Since(start))
}
