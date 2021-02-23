package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("./templates/*.gohtml"))

var fileserver = http.FileServer(http.Dir("./resources"))

func main() {
	fmt.Println("Server if started and listening on port 8080")
	http.HandleFunc("/", rootHandler)
	http.Handle("/resources/", http.StripPrefix("/resources", fileserver))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("Requested: %s -> %s\n", req.Method, req.URL.Path)

	w.Header().Set("Content-Type", "text/html")

	switch {
	case strings.HasSuffix(req.URL.Path, ".html"):
		f := req.URL.Path[1:len(req.URL.Path)-5] + ".gohtml"
		if err := templates.ExecuteTemplate(w, f, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		fmt.Println("File not found:", req.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}
}
