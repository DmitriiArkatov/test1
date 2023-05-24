package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// main - многопоточный генератор уникальных чисел
func main() {
	ch := make(chan int)
	trechmap := make(map[int]int)
	var uniquenum []int
	start := time.Now()
	var wg sync.WaitGroup

	limit, flow, count := oPtions()

	for i := 0; i < flow; i++ {
		wg.Add(1)
		go randomaizer(&wg, ch, limit, count)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for r := range ch {
		if r != trechmap[r] {
			fmt.Println(r, trechmap[r])
			uniquenum = append(uniquenum, r)
			trechmap[r] = r
			fmt.Println(uniquenum)
		}
		if len(uniquenum) == count {
			break
		}
	}

	fmt.Println(uniquenum)
	duration := time.Since(start)
	fmt.Printf("Время исполнения: %s\n", duration)
}

// Ввод аргументов с консоли
func oPtions() (int, int, int) {
	limit := flag.Int("limit", 15, "Limiting number generation")
	flow := flag.Int("flow", 5, "number of threads")
	count := flag.Int("count", 10, "number of unique numbers")
	flag.Parse()
	return *limit, *flow, *count
}

func randomaizer(wg *sync.WaitGroup, ch chan int, limit int, count int) {
	defer wg.Done()
	for i := 0; i <= count; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
