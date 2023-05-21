package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

var links []string
var text string //выжимка с сайтовw
// var links []string //срез для ссылок
func arguMents(argORC *string, argCreate *string) {
	fmt.Print("Введите путь до файла :")
	fmt.Scan(argORC)
	fmt.Print("Введите путь до директории:")
	fmt.Scan(argCreate)
} //считывает ввод и сохраняет в переменных

func fileORC(argORC *string) error {
	file, err := os.Open(*argORC) //путь до файла
	if err != nil {
		return err
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
		return err
	}
	fmt.Println(len(links))
	fmt.Println(links)
	return err
} //открываем файл , читаем и заносим в срез для ссылок, закрываем (file Open, Read, Close)

func scrap(i int) (string, error) {
	resp, err := http.Get(links[i])
	if err != nil {
		return "Fail", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	text = string(body)
	return "Done", err
} //забираем HTML из странички

func create(i string, argCreate *string) (string, error) {
	f, err := os.Create(*argCreate + "/сайт" + i) //путь до директории
	if err != nil {
		return "Fail", err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		return "Fail", err
	}
	return "Done", err
} //создаем новый файл и записываем в него то что забрали

/*
для проверки:
Введите путь до файла :
Введите путь до директории:
*/
func main() {
	argORC := new(string)        // переменная для считывания файла
	argCreate := new(string)     // переменная для создания новых файлов в директории
	arguMents(argORC, argCreate) // считывает ввод
	err := fileORC(argORC)       // открываем файл , читаем , закрываем
	if err != nil {
		fmt.Println(err)
	}
	for i := range links {
		_, err := scrap(i) //парсим страничку из файл
		if err != nil {
			fmt.Println(err)
		}
		i += 1
		i := strconv.Itoa(i)
		_, err = create(i, argCreate) //создаем новый файл в директории
		if err != nil {
			fmt.Println(err)
		} //достаем из среза ссылку , забираем HTML и создаем для него файл
	}
}
