package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	//получение исходных данных
	root := flag.String("root", ".", "Root directory to start traversal")
	sizeLimit := flag.Int64("size", 0, "Size limit in bytes")
	flag.Parse()
	//обходим дирректорию
	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//проверка являится ли директорией
		if info.IsDir() {
			dirSize, err := getDirSize(path)
			if err != nil {
				return err
			}
      //сравниваем размер директории с заданным параметром
			if dirSize > *sizeLimit {
				if err == writeToFile(fmt.Sprintf("%s: %d bytes\n", path, dirSize)) {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

// узнаем размер директории
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// добавляем в файл
func writeToFile(data string) error {
	file, err := os.OpenFile("infoDir", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	//проверка на то закрыли ли мы файл
	_, err = file.WriteString(data)
	return err
}
