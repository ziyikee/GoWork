package dir1

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

func test1() {
	// ruleid: waitgroup-add-called-inside-goroutine
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var x int32 = 0
	wg1.Add(1)
	for i := 0; i < 100; i++ {
		go func() {
			wg1.Add(1)
			go func() {
				wg2.Add(1)
				return
			}()
			atomic.AddInt32(&x, 1)
			wg1.Done()
			wg1.Add(1)
		}()
	}

	fmt.Println("Wait ...")
	wg1.Wait()
	fmt.Println(atomic.LoadInt32(&x))
}
func test1_3(errorCh chan error) {

	go func() {
		errorCh <- errors.New("aaa")
	}()
}
