package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	ch := make(chan string)
	filepath, dirpath := arguMent()      //считываем путь до файла(orc),который будем читать, и путь до директории для записи(pc)
	err, links := openingAfile(filepath) //забираем ошибку и срез с ссылками
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range *links {
		go parsCreate(i, *links, dirpath, ch) //парсим и записываем
		fmt.Printf("Сайт %s готов \n", <-ch)
	}
}

func arguMent() (*string, *string) {
	filepath := flag.String("pathfile", " ", "the path to the text file to be scanned ")          // переменная для считывания файла
	dirpath := flag.String("pathdir", " ", "the path to the directory for creating page content") // переменная для создания новых файлов в директории
	flag.Parse()
	return filepath, dirpath
} //считывает ввод и сохраняет в переменных

func openingAfile(argORC *string) (error, *[]string) {
	var links []string
	file, err := os.Open(*argORC) //путь до файла
	if err != nil {
		return err, nil
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
		return err, nil
	}
	fmt.Println(len(links))
	fmt.Println(links)
	return err, &links
} //открываем файл , читаем и заносим в срез для ссылок, закрываем (file Open, Read, Close)

func parsCreate(i int, links []string, pc *string, ch chan string) {
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
	f, err := os.Create(*pc + "/сайт" + n) //путь до директории
	if err != nil {
		fmt.Println(err)
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
	ch <- n

} //парсим страницу , создаем файл и записываем
