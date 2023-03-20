package dir1

import (
	"fmt"
	"sync"
)

func main2() {
	test2()
}

func test2() {
	// ruleid: waitgroup-add-called-inside-goroutine
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	wg3 := sync.WaitGroup{}
	wg3.Wait()
	wg1.Add(1)
	var list = []int{1, 2, 3, 4, 5}
	for i := 0; i < 100; i++ {
		go func() {
			wg1.Add(1)
			go func() {
				wg1.Add(1)
				go func() {
					return
				}()
				return
			}()
			wg1.Done()
			wg2.Add(1)
			addCall(wg2)
		}()
	}

	for i := 0; i < 100; i++ {

	}

	for _, i := range list {
		i++
	}

	fmt.Println("Wait ...")
	wg1.Wait()
}

func addCall(wg2 sync.WaitGroup) {
	wg2.Add(1)
}
