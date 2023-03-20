package dir1

import (
	"fmt"
	"time"
)

var (
	result = ""
)

func req1(timeout time.Duration) string {
	ch1 := make(chan string)
	// ruleid: hanging-goroutine
	go func() {
		newData := test()
		ch1 <- newData // block
	}()
	select {
	case result = <-ch1:
		fmt.Println("case result")
		return result
	case <-time.After(timeout):
		fmt.Println("case time.Afer")
		return ""
	}
}

func req1_FP(timeout time.Duration) string {
	ch3 := make(chan string, 1)
	// ok: hanging-goroutine
	go func() {
		newData := test()
		ch3 <- newData // block
	}()
	select {
	case result = <-ch3:
		fmt.Println("case result")
		return result
	case <-time.After(timeout):
		fmt.Println("case time.Afer")
		return ""
	}
}

func req2(timeout time.Duration) string {
	ch2 := make(chan string)
	// ruleid: hanging-goroutine
	go func() {
		newData := test()
		ch2 <- newData // block
	}()
	select {
	case <-ch2:
		fmt.Println("case result")
		return result
	case <-time.After(timeout):
		fmt.Println("case time.Afer")
		return ""
	}
}

func req2_FP(timeout time.Duration) string {
	ch4 := make(chan string, 1)
	// ok: hanging-goroutine
	go func() {
		newData := test()
		ch4 <- newData // block
	}()
	select {
	case <-ch4:
		fmt.Println("case result")
		return result
	case <-time.After(timeout):
		fmt.Println("case time.Afer")
		return ""
	}
}

func test() string {
	time.Sleep(time.Second * 2)
	return "very important data"
}
