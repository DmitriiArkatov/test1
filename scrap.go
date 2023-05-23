package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	filepath, dirpath := arguMents()
	if filepath == nil || dirpath == nil {
		panic("Wrong path!!")
		return
	}
	links, err := readSourse(filepath)
	if err != nil {
		return
	}
	for _, link := range links {
		writTodir(link, dirpath) //парсим и записываем

	}
	duration := time.Since(start)
	fmt.Printf("Время исполнения: %s\n", duration)
}

func arguMents() (*string, *string) {
	filepath := flag.String("pathfile", " ", "the path to the text file to be scanned ")          // переменная для считывания файла
	dirpath := flag.String("pathdir", " ", "the path to the directory for creating page content") // переменная для создания новых файлов в директории
	flag.Parse()
	return filepath, dirpath
}

func readSourse(filepath *string) ([]string, error) {
	var links []string
	file, err := os.Open(*filepath) //путь до файла
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
	scanner := bufio.NewScanner(file) //возвращает каждую строку текста, очищенную от маркеров конца строки. Возвращаемая строка может быть пустой. Маркером конца строки является один необязательный возврат каретки, за которым следует одна обязательная новая строка.
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() == true {
		links = append(links, scanner.Text())
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	fmt.Println(len(links))
	fmt.Println(links)
	return links, err
} //открываем файл , читаем и заносим в срез для ссылок, закрываем (file Open, Read, Close)

func writTodir(link string, dirpath *string) {
	//забираем страницу
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
	//тут происходит создание файла
	var i int
	i += 1
	n := strconv.Itoa(i)                        // преобразуем в строку
	f, err := os.Create(*dirpath + "/сайт" + n) //путь до директории
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(f)
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}

} //парсим страницу , создаем файл и записываем
