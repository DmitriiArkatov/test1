package main

import (
	"github.com/gorilla/websocket"
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
func createLocalhost() {
	http.HandleFunc("/", local)
	http.HandleFunc("/ws", wsLocal)
}
func local(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/vitaliy/awesomeProject/trening/response.html")
}
func wsLocal(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	log.Println("Client Connected")

	reader(conn)
}
func reader(conn *websocket.Conn) {
	var (
		limit, flow, count int
	)
	ch := make(chan int)
	defer close(ch)

	go sendElements(conn, ch)
	for {
		// прочитать в сообщении
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := strings.Split(string(p), ",")
		// выводим это сообщение для наглядности
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
		go Random(limit, flow, count, ch)
	}
}
func sendElements(conn *websocket.Conn, ch <-chan int) {
	for {
		r, ok := <-ch
		if !ok {
			// Канал закрыт, передача элементов завершена
			return
		}

		// Преобразуем элемент в строку и отправляем клиенту
		answer := strconv.Itoa(r)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
			log.Println(err)
			return
		}
		time.Sleep(time.Second * 2)
	}
}
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
func randomaizer(wg *sync.WaitGroup, ch chan int, limit int) {
	defer wg.Done()
	for i := -1; i <= 100; i++ {
		r := rand.Intn(limit + 1)
		ch <- r
	}
}
