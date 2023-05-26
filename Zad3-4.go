package Zadacha3

import (

	"math/rand"
	"sync"

)

// main - многопоточный генератор уникальных чисел
func Random(limit int, flow int, count int) []int {
	ch := make(chan int)
	trashmap := make(map[int]int)
	var uniquenum []int
	var wg sync.WaitGroup
	for i := 0; i < flow; i++ {
		wg.Add(1)
		go Randomaizer(&wg, ch, limit)
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
	return uniquenum
}


// генератор уникальных чисел
func Randomaizer(wg *sync.WaitGroup, ch chan int, limit int) {
	defer wg.Done()
	for i := -1; i <= 100; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
