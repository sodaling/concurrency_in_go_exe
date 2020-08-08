package main

import "time"

func main() {
	// 这边弄了个自动产生心跳，还会每两个心跳周期产生一个结果的函数
	doWork := func(done <-chan struct{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		// 注意，返回的heartbeat是interface类型的
		heartbeat := make(chan interface{})
		result := make(chan time.Time)
		go func() {
			defer close(result)
			defer close(heartbeat)
			// 一个是产生结果的，第二个是心跳
			// 两个返回都是当前的时间
			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)
			// 首先是发送心跳的，这里得注意的是，发送心跳不能堵塞。因为不一定每个调用的人都care心跳

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
					//记住千万不能堵塞,所以也不用担心done
				}
			}
			sendResult := func(item time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						// 还是得发送心跳
						sendPulse()
					case result <- item:
						return
					}
				}
			}
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case item := <-workGen:
					sendResult(item)
				}
			}
		}()
		// 重点就是
		// 有堵塞的地方，就记得关注done和发送心跳
		return heartbeat, result
	}
}
