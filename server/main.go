package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", HelloHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error during server start: %v", err)
	}
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "Hello Go")
	if err != nil {
		log.Fatalf("Error while writing response")
	}
}
