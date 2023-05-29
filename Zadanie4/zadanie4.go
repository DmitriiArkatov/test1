package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type FormData struct {
	Limit string `json:"limit"`
	Flow  string `json:"flow"`
	Count string `json:"count"`
}

func main() {
	createLocalhost()
}

func local(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "/home/vitaliy/awesomeProject/trening/primer.html")
	}
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Неверный Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		// Распарсим JSON-тело запроса в структуру FormData
		var formData FormData
		err := json.NewDecoder(r.Body).Decode(&formData)
		if err != nil {
			http.Error(w, "Ошибка при чтении JSON-данных", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
		}(r.Body)
		// Делаем что-то с данными
		limit, err := strconv.Atoi(formData.Limit)
		if err != nil {
			http.Error(w, "Ошибка при конвертировании", http.StatusBadRequest)
		}
		flow, err := strconv.Atoi(formData.Flow)
		if err != nil {
			http.Error(w, "Ошибка при конвертировании", http.StatusBadRequest)
		}
		count, err := strconv.Atoi(formData.Count)
		if err != nil {
			http.Error(w, "Ошибка при конвертировании", http.StatusBadRequest)
		}
		uniquenum := Random(limit, flow, count)
		fmt.Println("Получено сообщение:", uniquenum)

		// Отправляем ответ
		response := uniquenum
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func createLocalhost() {
	port := "8888"
	http.HandleFunc("/", local)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Random - многопоточный генератор уникальных чисел
func Random(limit int, flow int, count int) []int {
	ch := make(chan int)
	trashmap := make(map[int]int)
	var uniquenum []int
	start := time.Now()
	var wg sync.WaitGroup

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
	return uniquenum
}

// генератор уникальных чисел
func randomaizer(wg *sync.WaitGroup, ch chan int, limit int) {
	defer wg.Done()
	for i := -1; i <= 100; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
