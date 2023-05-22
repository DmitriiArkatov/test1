package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func main() {
	ch := make(chan string)

	ors, pc := arguMents()     //считываем путь до файла(orc),который будем читать, и путь до директории для записи(pc)
	err, links := fileORC(ors) //забираем ошибку и срез с ссылками
	if err != nil {
		fmt.Println(err)
	}
	for i := range *links {
		go parsCreate(i, *links, pc, ch) //парсим и записываем
		fmt.Println(<-ch)
	}
}

func arguMents() (*string, *string) {
	argORC := new(string)    // переменная для считывания файла
	argCreate := new(string) // переменная для создания новых файлов в директории
	fmt.Print("Введите путь до файла :")
	fmt.Scan(argORC)
	fmt.Print("Введите путь до директории:")
	fmt.Scan(argCreate)
	return argORC, argCreate
} //считывает ввод и сохраняет в переменных

func fileORC(argORC *string) (error, *[]string) {
	var links []string
	file, err := os.Open(*argORC) //путь до файла
	if err != nil {
		return err, nil
	}
	defer file.Close()
	//var str []string
	scanner := bufio.NewScanner(file) //возвращает каждую строку текста, очищенную от маркеров конца строки. Возвращаемая строка может быть пустой. Маркером конца строки является один необязательный возврат каретки, за которым следует одна обязательная новая строка.
	scanner.Split(bufio.ScanLines)
	//var links *[]string
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
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	text := string(body)

	//тут происходит создание файла
	i += 1
	n := strconv.Itoa(i)
	f, err := os.Create(*pc + "/сайт" + n) //путь до директории
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		fmt.Println(err)
	}
	ch <- "Сайт " + n + " готов"

} //парсим страницу , создаем файл и записываем
