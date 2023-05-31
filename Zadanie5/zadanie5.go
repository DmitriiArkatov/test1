package main

import (
	"fmt"
	"gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	createLocalhost()
	log.Fatal(http.ListenAndServe(":8888", nil))
}

// createLocalhost - создает локальные сервера
func createLocalhost() {
	http.HandleFunc("/", local)
	http.HandleFunc("/ws", wsLocal)
}

// local - сервер с соединение HTTP
func local(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/vitaliy/awesomeProject/trening/response.html")
}

// wslocal - сервер с соединение WebSocket
func wsLocal(w http.ResponseWriter, r *http.Request) {
	//буфер чтения и записи для WebSocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //Проверяем входящий запрос на подключение по WebSocket

	conn, err := upgrader.Upgrade(w, r, nil) //Модернизируем наше HTTP соединение и подключаемся по WebSocket
	if err != nil {
		log.Println(err)
		return
	}

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	log.Println("Client connected")

	reader(conn)
}

// reader - читаем сообщение по WebSocket
func reader(conn *websocket.Conn) {
	wg := sync.WaitGroup{}

	var (
		limit, flow, count int
	)
	ch := make(chan int)
	defer close(ch)

	for {
		_, p, err := conn.ReadMessage() // читаем сообщения
		if err != nil {
			log.Println(err)
			return
		}
		message := strings.Split(string(p), ",") //конвертируем полученое сообщение в массив
		fmt.Println(message)
		//парсим сообщение-массив
		for n, i := range message {
			n++
			switch n {
			case 1:
				limit, _ = strconv.Atoi(i)
			case 2:
				flow, _ = strconv.Atoi(i)
			case 3:
				count, _ = strconv.Atoi(i)
			}
		}
		fmt.Println(limit, flow, count)
		//if limit > 0 && flow > 0 && count > 0 {

		if limit > 0 && flow > 0 && count > 0 {
			wg.Add(1)
			go Random(limit, flow, count, ch)
			go sendAnswer(conn, ch, &wg)
			
		} else {
			answer := "Все параметры должны быть > 0!!!"
			if err := conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// sendAnswer - отсылаем ответ пользователю по WebSocket
func sendAnswer(conn *websocket.Conn, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Wait()
	for {
		r, ok := <-ch
		if !ok {
			return
		}

		answer := strconv.Itoa(r) // Преобразуем элемент в строку и отправляем клиенту
		//Отправка ответа пользователю
		if err := conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
			log.Println(err)
			return
		}
		time.Sleep(time.Second)
	}
}

// Random - получаем уникальные числа по заданным параметрам
func Random(limit int, flow int, count int, ch1 chan<- int) {
	ch := make(chan int)
	trashmap := make(map[int]int)
	var uniquenum []int
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
			ch1 <- r
		}
		if len(uniquenum) == count {
			break
		}
	}
}

// randomaizer - генерирует рандомные числа
func randomaizer(wg *sync.WaitGroup, ch chan int, limit int) {
	defer wg.Done()
	for i := -1; i <= 100; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
