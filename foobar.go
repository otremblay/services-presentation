package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	log.Println("Mapping handlers...") //OMIT
	http.Handle("/foo", fooHandler)

	http.HandleFunc("/bar/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Println("Starting server...") //OMIT
	log.Fatal(http.ListenAndServe(":12345", nil))
	log.Println("Stopped server") //OMIT
}

var fooHandler http.HandlerFunc = func(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(rw, "baz")
}
