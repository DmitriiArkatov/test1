package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// main - многопоточный генератор уникальных чисел
func main() {
	ch := make(chan int)
	trashmap := make(map[int]int)
	var uniquenum []int
	start := time.Now()
	var wg sync.WaitGroup

	limit, flow, count := options()

	for i := 0; i < flow; i++ {
		wg.Add(1)
		go randomaizer(&wg, ch, limit)
	}
	//запускаем горутину которая будет закрывать каналы и остальные горутины
	go func() {
		wg.Wait()
		close(ch)
	}()
	//проверка уникальности числа в списке
	for r := range ch {
		if r != trashmap[r] {
			uniquenum = append(uniquenum, r)
			trashmap[r] = r
		}
		if len(uniquenum) == count {
			break
		}
	}
	//вывод списка уникальных чисел и  время исполнения
	fmt.Println(uniquenum)
	duration := time.Since(start)
	fmt.Printf("Время исполнения: %s\n", duration)
}

// Ввод аргументов с консоли
func options() (int, int, int) {
	limit := flag.Int("limit", 10, "Limiting number generation") //до какого числа генерировать
	flow := flag.Int("flow", 1, "number of threads")             //кол - во горутин
	count := flag.Int("count", 10, "number of unique numbers")   //кол - во уникальных элементов в списке
	flag.Parse()
	if *limit == 0 && *count == 0 {
		log.Fatal("limit and count must be > 0")
	}
	return *limit, *flow, *count
}

// генератор уникальных чисел
func randomaizer(wg *sync.WaitGroup, ch chan int, limit int) {
	defer wg.Done()
	for i := -1; i <= 100; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
