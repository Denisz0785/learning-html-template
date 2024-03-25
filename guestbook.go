package main

import (
	"bufio"
	"fmt"
	"html/template"
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
	html, err := template.ParseFiles("view.html")
	check(err)
	guestbook := Guestbook{
		SignatureCount: len(signatures),
		Signatures:     signatures,
	}
	err = html.Execute(w, guestbook)
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

type Guestbook struct {
	SignatureCount int
	Signatures     []string
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	html, err := template.ParseFiles("new.html")
	check(err)
	err = html.Execute(w, nil)
	check(err)
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
func creatHandler(w http.ResponseWriter, r *http.Request) {
	sig := r.FormValue("signature")
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE
	file, err := os.OpenFile("signature.txt", options, os.FileMode(0600))
	check(err)
	_, err = fmt.Fprintln(file, sig)
	check(err)
	err = file.Close()
	check(err)
	http.Redirect(w, r, "/guestbook", http.StatusFound)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/guestbook", viewHandler)
	mux.HandleFunc("/guestbook/new", newHandler)
	mux.HandleFunc("/guestbook/create", creatHandler)
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}
