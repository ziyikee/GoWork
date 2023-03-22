package dir1

import (
	"fmt"
	"testing"
	"time"
)

var values = []int{1, 2, 3, 4, 5}

func Testa(t *testing.T) {
	for _, value := range values {
		go func() {
			fmt.Printf("value = %d\n", value)
		}()
	}

	time.Sleep(1 * time.Second)
	fmt.Println("")

	for _, value := range values {
		go func(value int) {
			fmt.Printf("value = %d\n", value)
		}(value)

	}

	time.Sleep(1 * time.Second)
	fmt.Println("")

	for _, value := range values {
		go func(a int) {
			fmt.Printf("value = %d\n", a)
		}(value)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("")
	var i = 0
	for ; i < len(values); i++ {
		go func() {
			if i > 0 {
				fmt.Printf("value = %d\n", values[i])
			}
			go func() {
				fmt.Printf("value = %d\n", values[i])
			}()
			var j = 0
			for ; j < len(values); j++ {
				if i > 0 {
					fmt.Printf("value = %d\n", values[j])
				}
				go func() {
					fmt.Printf("value = %d\n", values[j])
				}()
			}
		}() //panic,values[5]
	}

	time.Sleep(1 * time.Second)
	fmt.Println("")

	for ; i < len(values); i++ {
		go func(a int) {
			fmt.Printf("value = %d\n", values[a])
		}(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("")

	for i := 0; i < len(values); i++ {
		go func(i int) {
			fmt.Printf("value = %d\n", values[i])
		}(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("")

	//16

	for ; i < len(values); i++ {

	}

	for i := 0; 0 < i; i-- {

	}

	for i := 0; i < len(values); i = i + 1 {

	}

	for i := 0; i > 0; i = i - 1 {

	}

	for i := 0; i < len(values); i += 1 {

	}

	for i := 0; i > 0; i -= 1 {

	}

	for i := 1; i < len(values); i = i * 2 {

	}

	for i := 1; i > 0; i = i / 2 {

	}

}
