package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// main - многопоточно парсит по этим ссылкам HTML-страницы  и создает локальные серверы
func main() {
	start := time.Now()
	var wg sync.WaitGroup
	filepath := arguments()
	links, err := readSources(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, link := range links {
		i += 1234
		wg.Add(1)
		go parser(link, &wg, i) //каждая горутина забирает HTML со страницы и отправляет его на наш сервер
		fmt.Printf("http://localhost:%d\n", i)
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Время исполнения: %s\n", duration)
}

// arguments - считываем аргументы с консоли
func arguments() string {
	filepath := flag.String("pathfile", "/home/vitaliy/awesomeProject/zadone/link.txt", "the path to the text file to be scanned ") // переменная для считывания файла
	flag.Parse()
	return *filepath
}

// readSources - читает файл и заносит все в срез
func readSources(argORC string) ([]string, error) {
	var links []string
	file, err := os.Open(argORC) //путь до файла
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	fmt.Println(len(links))
	fmt.Println(links)
	return links, nil
}

// получаем HTML страницы
func parser(link string, wg *sync.WaitGroup, i int) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	text := string(body)
	createLocalhost(text, i, wg)
}

// создание сервера на который отправляем HTML
func createLocalhost(text string, por int, wg *sync.WaitGroup) {
	defer wg.Done()
	port := strconv.Itoa(por)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, text)
	})

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}
