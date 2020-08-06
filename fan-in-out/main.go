package main

import (
	"fmt"
	"sync"
)

func main() {
	// 就好像一个班级的作业，原本由1位老师去修改，现在变成了8位老师同时批改
	// 可以考虑作为爬虫的架构部分。从很多个去爬，然后写入的话集中一个地方（比如文本）
	// 但是其实一般都是写入数据库而不是文本

	//扇入意味着将多个数据流复用或合并成一个流。
	// 注意，对返回结果有顺序要求的情况下不适合这个模式，因为我们没法保证那么多个channel每个的到达顺序
	fanIn := func(done <-chan interface{}, channels ...<-chan int) <-chan interface{} {
		ret := make(chan interface{})
		// 这是保证所有通道组件都要读取完
		var wg sync.WaitGroup
		multiplex := func(c <-chan int) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case ret <- <-c:
				}
			}
		}

		//添加需要等到的数
		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}
		go func() {
			// 其实简单来说也就是启了多个go程序往返回的go程序里塞
			// 启了个go程来快速返回，一般不这么做，但是这里这么做贼合适
			wg.Wait()
			//全部读完了，记得关闭
			close(ret)
		}()

		return ret
	}
	//这边就不实验了，因为书上的一个求素数的函数没写全
	// 书上的效果不错，时间降低了百分之78

	fmt.Println(fanIn)
}
