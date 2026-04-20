package main

import (
	"fmt"
	"sync"
)

// print 1 to 10 in squance
// by the help of two gorotuine where one generate even number and second one odd number

func main1() {
	ch := make(chan int)
	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	wg.Add(2)
	go func() { // even
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			m.Lock()
			ch <- i
			m.Unlock()

		}

	}()

	go func() { // odd
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			m.Lock()
			ch <- i
			m.Unlock()
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println("got", v)

	}

}
