package main

import (
	"flag"
	"fmt"
	"math/rand"
)

func main() {
	//var wg sync.WaitGroup
	var uniquenum []int
	limit, flow, count := options()
	fmt.Println("началось")
	for len(uniquenum) != count {
		for i := 0; i < flow; i++ {
			//wg.Add(1)
			fmt.Println("началось")
			filtr(&limit, uniquenum /*&wg*/)
		}
		//wg.Wait()
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

func filtr(limit *int, uniquenum []int /*, wg *sync.WaitGroup*/) {
	r := rand.Intn(*limit)
	bo := Contains(uniquenum, r)
	if bo != true {
		uniquenum = append(uniquenum, r)
	} else {
	}
	//wg.Done()
}

func Contains(uniquenum []int, x int) bool {
	for _, n := range uniquenum {
		if x == n {
			return true
		}
	}
	return false
}
