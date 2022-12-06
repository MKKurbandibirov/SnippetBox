package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippets)
	mux.HandleFunc("/snippet/create", createSnippet)

	log.Println("Start!")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
