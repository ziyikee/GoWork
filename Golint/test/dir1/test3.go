package dir1

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	test3_1()
	test3_2()
	test3_3()
}

func test3_1() {
	// ruleid: waitgroup-wait-inside-loop
	wg1 := sync.WaitGroup{}
	var x int32 = 0
	for i := 0; i < 100; i++ {
		wg1.Add(1)
		go func() {
			wg1.Done()
			atomic.AddInt32(&x, 1)
		}()
		wg1.Wait()
	}

	fmt.Println("Wait ...")
	fmt.Println(atomic.LoadInt32(&x))
}

func test3_2() {
	// ruleid: waitgroup-wait-inside-loop
	var wg2 sync.WaitGroup
	var x int32 = 0
	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func() {
			atomic.AddInt32(&x, 1)
			wg2.Done()
		}()
		wg2.Wait()
	}

	fmt.Println("Wait ...")
	fmt.Println(atomic.LoadInt32(&x))
}

func test3_3() {
	// ok: waitgroup-wait-inside-loop
	var wg3 sync.WaitGroup
	var x int32 = 0
	for i := 0; i < 100; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomic.AddInt32(&x, 1)
		}()
	}

	fmt.Println("Wait ...")
	wg3.Wait()
	fmt.Println(atomic.LoadInt32(&x))
}
func test3_4() {
	// ok: waitgroup-wait-inside-loop

	var x int32 = 0
	for i := 0; i < 100; i++ {
		var wg4 sync.WaitGroup
		wg4.Add(1)
		go func() {
			defer wg4.Done()
			atomic.AddInt32(&x, 1)
		}()
		wg4.Wait()
	}

	fmt.Println("Wait ...")
	fmt.Println(atomic.LoadInt32(&x))
}
