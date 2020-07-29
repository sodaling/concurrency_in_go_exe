package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

func main() {
	// 尽量减少锁的范围，同时defer来做解锁操作防止panic时候解锁不成功
	// n个生产者，每个生产者拿锁5次，每次间隔1nano秒
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			l.Lock()
			l.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		// 在这里不区分rw或者一般的locker
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		// 启动一个生产者
		go producer(&wg, mutex)
		for i := 0; i < count; i++ {
			// 这边用读写锁
			go observer(&wg, rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()
	var m sync.RWMutex

	fmt.Fprintf(tw, "Readers\tRwMutext\tmutex\n")
	for i := 0; i < 20; i++ {
		// 20次方
		// 注意RLocker会返回一个实现了locker接口的读锁
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw, "%d\t%v\t%v\n", count, test(count, &m, m.RLocker()), test(count, &m, &m),
		)
	}
}
