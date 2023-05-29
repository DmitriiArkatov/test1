package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type FormData struct {
	Limit int `json:"limit"`
	Flow  int `json:"flow"`
	Count int `json:"count"`
}

func main() {
	createLocalhost()
}

func local(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "primerno.html")
	case http.MethodPost:
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

		// Делаем что-то с данными
		fmt.Println("Получено сообщение:", formData.Flow, formData.Count, formData.Limit)

		// Отправляем ответ
		response := struct {
			Status string `json:"status"`
		}{
			Status: "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	}
	//if r.Method != http.MethodPost {
	//	http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
	//	return
	//}

	// Проверяем Content-Type заголовок

}

//func primer(w http.ResponseWriter, r *http.Request) {
//	t, err := template.ParseFiles("/home/vitaliy/awesomeProject/inputBox.html")
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//	}
//	err = t.ExecuteTemplate(w, "button", nil)
//	if err != nil {
//		fmt.Fprintf(w, err.Error())
//	}
//
//}
func createLocalhost() {
	port := "8888"
	mux := http.NewServeMux()

	//mux.HandleFunc("/", primer)
	mux.HandleFunc("/", local)

	err := http.ListenAndServe(":"+port, mux)
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
