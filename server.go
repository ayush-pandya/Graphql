package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/tickets", helloHandler)
	fmt.Printf("starting a server port:8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tickets" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	// Simulate a delay
	time.Sleep(2 * time.Second)
	// Return a simple JSON response
	fmt.Fprintf(w, `{"message": "Hello, World!"}`)
	// Return a simple text response
	// fmt.Fprintf(w, `{"message": "Hello, World!"}`)
	fmt.Fprintf(w, "Hello!")
}
