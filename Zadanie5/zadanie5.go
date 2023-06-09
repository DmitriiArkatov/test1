package main

import (
	"gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	createLocalhost()
	log.Fatal(http.ListenAndServe(":"+"8888", nil))
}

// createLocalhost - создает локальные сервера
func createLocalhost() {
	http.HandleFunc("/", local)
	http.HandleFunc("/front.js", localjs)
	http.HandleFunc("/ws", wsLocal)
}

// local - сервер с соединение HTTP
func local(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/vitaliy/awesomeProject/trening/zadanie5.html")
}

// localjs - выгружаем JS файл
func localjs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/vitaliy/awesomeProject/trening/front.js")
}

// wslocal - сервер с соединение WebSocket
func wsLocal(w http.ResponseWriter, r *http.Request) {
	//буфер чтения и записи для WebSocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//забираем из url данные
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	flow, _ := strconv.Atoi(r.URL.Query().Get("flow"))
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } //Проверяем входящий запрос на подключение по WebSocket

	conn, err := upgrader.Upgrade(w, r, nil) //Модернизируем наше HTTP соединение и подключаемся по WebSocket
	if err != nil {
		log.Println(err)
		return
	}

	//закрываем соединение
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	log.Println("Client connected")
	
	randomAnswer(conn, limit, flow, count)
}

// randomAnswer - создаем уникальный ответ перед отправкой
func randomAnswer(conn *websocket.Conn, limit int, flow int, count int) {
	ch := make(chan int)
	var wg sync.WaitGroup
	defer wg.Wait()

	//если условие - истина , то начинается запись уникального числа в канал, для отправки
	if limit > 0 && flow > 0 && count > 0 {
		go checkSendAnswer(conn, ch, count, &wg)
		for i := 0; i < flow; i++ {
			//запускается указаное кол-во горутин
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := -1; i <= 100; i++ {
					r := rand.Intn(limit + 1)
					ch <- r
				}
			}()
		}
	} else {
		answer := "Все параметры должны быть > 0!!!"
		if err := conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
			log.Println(err)
			return
		}
	}
}

// checkSendAnswer - проверяем на уникальность и отсылаем ответ пользователю по WebSocket
func checkSendAnswer(conn *websocket.Conn, ch chan int, count int, wg *sync.WaitGroup) {
	trashmap := make(map[int]int)
	defer close(ch)
	defer wg.Wait()
	//проверка уникальности числа в списке
	for r := range ch {
		if r != trashmap[r] {
			trashmap[r] = r
			answer := strconv.Itoa(r) // Преобразуем элемент в строку и отправляем клиенту
			//Отправка ответа пользователю
			if err := conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
				log.Println(err)
				return
			}
			time.Sleep(time.Second)
		}
		if len(trashmap) == count {
			break
		}
	}
}
