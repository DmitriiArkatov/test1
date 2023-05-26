package main

import (
	"fmt"
	"html/template"
	"net/http"
	"packet"
	"strconv"
)



func main() {
	createLocalhost()
}

func local(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("/home/vitaliy/awesomeProject/inputBox.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	err = t.ExecuteTemplate(w, "button", nil)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	if r.Method == "POST" {
		limit, _ := strconv.Atoi(fmt.Sprintf(r.FormValue("limit")))
		flow, _ := strconv.Atoi(fmt.Sprintf(r.FormValue("flow")))
		count, _ := strconv.Atoi(fmt.Sprintf(r.FormValue("count")))

		if limit == 0 || count == 0 {
			fmt.Fprintf(w, "Please enter the correct data")
		}
		uniquenum := Zadacha3.Random(limit, flow, count)
		fmt.Fprintf(w, "%d", uniquenum)

	}
}

func createLocalhost() {
	port := "8888"
	mux := http.NewServeMux()

	mux.HandleFunc("/", local)

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println(err)
		return
	}
}
