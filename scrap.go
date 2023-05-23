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
	duration := time.Since(start)
	filepath, dirpath := arguMent()
	if filepath == nil || dirpath == nil {
		panic("Wrong path!!")
		return
	}
	err, links := openingAfile(filepath)
	if err != nil {
		return
	}
	for i := range *links {
		parsCreate(i, *links, dirpath) //парсим и записываем

	}
	fmt.Printf("Время исполнения: %s\n", duration)
}

func arguMent() (*string, *string) {
	filepath := flag.String("pathfile", " ", "the path to the text file to be scanned ")          // переменная для считывания файла
	dirpath := flag.String("pathdir", " ", "the path to the directory for creating page content") // переменная для создания новых файлов в директории
	flag.Parse()
	return filepath, dirpath
}

func openingAfile(filepath *string) (error, *[]string) {
	var links []string
	file, err := os.Open(*filepath) //путь до файла
	if err != nil {
		return err, nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)
	scanner := bufio.NewScanner(file) //возвращает каждую строку текста, очищенную от маркеров конца строки. Возвращаемая строка может быть пустой. Маркером конца строки является один необязательный возврат каретки, за которым следует одна обязательная новая строка.
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() == true {
		links = append(links, scanner.Text())
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err, nil
	}
	fmt.Println(len(links))
	fmt.Println(links)
	return err, &links
} //открываем файл , читаем и заносим в срез для ссылок, закрываем (file Open, Read, Close)

func parsCreate(i int, links []string, dirpath *string) {
	//забираем страницу
	resp, err := http.Get(links[i])
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
	i += 1
	n := strconv.Itoa(i)
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
