package dir1

func example() {
	var ch1 = make(chan string)
	ch2 := make(chan string, 1)
	var ch3, ch5 = make(chan string), make(chan string, 0)
	ch4 := make(chan string, 0)

	go func() {
		str := "aaa"
		ch1 <- str
		ch2 <- str
		ch3 <- str
		ch4 <- str
		ch5 <- str
	}()
	a := 1
	select {
	case <-ch1:
		a++
	case b := <-ch2:
		a = len(b)
	case <-ch3:
		a++
	case <-ch5:
		a++
	}
}
