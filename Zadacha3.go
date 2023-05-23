package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	ch := make(chan int)
	var uniquenum []int
	limit, flow, count := options()
	fmt.Println("началось")
	for len(uniquenum) == count {
		for i := 0; i <= flow; i++ {
			fmt.Println("началось")
			go randomaizer(limit, ch)
		}
		for i := 0; i <= flow; i++ {
			fmt.Println(<-ch)
			filtr(ch, uniquenum)
		}
	}
	fmt.Println(uniquenum)
}

// Ввод аргументов с консоли
func options() (int, int, int) {
	limit := flag.Int("limit", 0, "Limiting number generation")
	flow := flag.Int("flow", 0, "number of threads")
	count := flag.Int("count", 0, "number of unique numbers")
	return *limit, *flow, *count
}

func randomaizer(limit int, ch chan int) {
	r := rand.Intn(limit)
	ch <- r

}

func filtr(ch chan int, uniquenum []int) {
	bo := Contains(uniquenum, <-ch)
	if bo != true {
		uniquenum = append(uniquenum, <-ch)
	} else {

	}
}

func Contains(uniquenum []int, x int) bool {
	for _, n := range uniquenum {
		if x == n {
			return true
		}
	}
	return false
}
