package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	signatures := getStrings("signature.txt")
	fmt.Printf("%v\n", signatures)

	html, err := template.ParseFiles("view.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)
	message := []byte("signature list goes here")
	_, err = w.Write(message)
	check(err)
	_, err = io.WriteString(w, signatures[0])
	check(err)
}
func executeTemplate(text string, data interface{}) {
	tmpl, err := template.New("test").Parse(text)
	check(err)
	err = tmpl.Execute(os.Stdout, data)
	check(err)

}

type Client struct {
	Name string
	Age  int
}

type Subscriber struct {
	Name   string
	Rate   float64
	Active bool
}

func getStrings(fileName string) []string {
	var lines []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	check(scanner.Err())
	return lines
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/guestbook", viewHandler)
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)

	text := "Action start\n {{.}} \nEnd"
	tmpl, err := template.New("test").Parse(text)
	check(err)
	err = tmpl.Execute(os.Stdout, 42)
	check(err)
	err = tmpl.Execute(os.Stdout, "om a hum")
	check(err)
	err = tmpl.Execute(os.Stdout, 34)
	check(err)

	executeTemplate("Hello {{.}} end\n", "tatata")
	executeTemplate("start {{if .}} true or false {{end}}finish\n", false)
	tmplTxt := "Before loop {{.}}\n {{range .}} in loop {{.}}\n {{end}} after loop{{.}}\n"
	executeTemplate(tmplTxt, []string{"do", "re", "mi"})
	tmplNum := "Prices: {{range .}} {{.}}\n {{end}}"
	executeTemplate(tmplNum, []float64{3.45, 5.67, 67.34})

	tmplStruct := "Data from structure {{.Name}}\n {{.Age}}\n"
	executeTemplate(tmplStruct, Client{"Denis", 38})

	tmplStr2 := "Subscriber {{.Name}}\n {{if .Active}} Rate {{.Rate}}\n {{end}}"
	sub := Subscriber{"Denis", 56.67, false}
	executeTemplate(tmplStr2, sub)
}
