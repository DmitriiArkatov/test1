package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// main - многопоточный генератор случайных чисел
func main() {
	var uniquenum []int
	start := time.Now()
	var wg sync.WaitGroup
	limit, flow, count := options()
	for len(uniquenum) != count {
		for i := 0; i <= flow; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := rand.Intn(limit + 1)
				bo := uniquenessCheck(uniquenum, r)
				if bo != true {
					uniquenum = append(uniquenum, r)
				} else {
				}
			}()
		}
		wg.Wait()
	}
	fmt.Println(uniquenum)
	duration := time.Since(start)
	fmt.Printf("Время исполнения: %s\n", duration)
}

// Ввод аргументов с консоли
func options() (int, int, int) {
	limit := flag.Int("limit", 0, "Limiting number generation")
	flow := flag.Int("flow", 0, "number of threads")
	count := flag.Int("count", 0, "number of unique numbers")
	flag.Parse()
	return *limit, *flow, *count
}

// uniquenessCheck - проверяет наличие числа в списке
func uniquenessCheck(uniquenum []int, x int) bool {
	for _, n := range uniquenum {
		if x == n {
			return true
		}
	}
	return false
}
