package main

import "fmt"

func main() {
	//将一些了通道拆解为一个简单的通道——我们成为通道的桥接，这使得消费者更加容易关注手头的问题
	// 首先还是需要orDone函数
	orDone := func(done <-chan struct{},inStream <-chan interface{}) <-chan interface{}{
		ret := make(chan interface{})
		go func() {
			defer  close(ret)
				for {
					select {
					case <-done:
						return
					case v,ok:= <-inStream:
						if ok==false{
							return
						}
						select {
						case <-done:
						case ret<-v:
						}

					}
				}
		}()
		return ret
	}

	//简单来说就是多个<-channel组合成一个channel
	bridge:= func(done<-chan struct{},chanStream<-chan <-chan interface{})<-chan interface {}{
		ret := make(chan interface{})
		go func() {
			defer close(ret)
			for {
				var stream <-chan interface{}
				select {
				// 先取出来
				case maybeStream,ok:=<-chanStream:
					if ok == false{
						return
					}
					stream = maybeStream

				case <-done:
					return
				}
				//然后把取出来的遍历
				for val := range orDone(done,stream){
					select {
					case ret <-val:
					case <-done:
					}
				}
			}
		}()
		return ret
	}
	// 首先需要一个生产一系列chan的chan
	genVal := func()<-chan <-chan interface{}{
		ret := make(chan (<-chan interface{}))
		go func() {
			defer close(ret)
			for i:=0;i<10;i++{
				stream := make(chan interface{},1)
				stream<-i
				close(stream)
				// 反正就构造一个只有一个值的chan，然后塞进ret返回
				ret<-stream
			}
		}()
		return ret
	}
	for v:=range bridge(nil,genVal()){
		// 然后我们看看效果
		fmt.Println(v)
	}
}
