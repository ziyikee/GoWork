package dir1

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	//functions := func(a int) {
	//	fmt.Printf("func: %d\n", a)
	//}

	for i := 0; i <= 10; i++ {
		go func(a int) {
			fmt.Printf("func: %d\n", i)
		}(i)
	}
	time.Sleep(3 * time.Second)
}
