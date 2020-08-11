package main

import (
	"fmt"
	"log"
	"time"
)

type startGoroutineFn func(done <-chan interface{}, pluseInterval time.Duration) (heartbeat <-chan interface{})

func main() {
	// 这边弄了个自动产生心跳，还会每两个心跳周期产生一个结果的函数
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
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

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {
		// 自动关闭
		close(done)
	})
	const timeout = 2 * time.Second

	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}

	}

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

	newSteward := func(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn {
		// 这边弄一个检测心跳，并且自动重新拉起来的函数
		return func(done <-chan interface{}, pulseInterval time.Duration) <-chan interface{} {
			heartbeat := make(chan interface{})
			go func() {
				defer close(heartbeat)
				// 下面两个ward开头的是专门给所监控的函数用的
				// 在外面声明就开在外面一直持有
				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}
				startWard := func() {
					wardDone = make(chan interface{})
					// 值得注意的是，这个函数其实并不返回结果，只返回一个心跳的chan
					// 监控程序心跳间隔是被监控时长的两倍。也就是里面包两次平安了，我外面才报平安一次
					wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
				}
				startWard()
				pulse := time.Tick(pulseInterval)
			monitorLoop:
				for {
					// 超时信号一旦过了，就说明这个已经不健康了，要重新拉起来了
					timeoutSignal := time.After(timeout)
					for {
						select {
						case <-pulse:
							// 监控程序的心跳得往外传送
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case <-wardHeartbeat:
							// 被监控程序的心跳得接受
							continue monitorLoop
						case <-timeoutSignal:
							// 被监控的心跳半天不接受，自己半天不来,那就得重启了
							log.Printf("steward: ward unhealthy;restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							// 外面叫停了
							return
						}
					}
				}
			}()
			return heartbeat
		}
	}

	//上面的是比较简单的，除了取消操作和心跳所需信息之外不接受也不返回任何结果。但是可以用闭包强化
}
