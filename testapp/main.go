package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/books", list)
	http.HandleFunc("/books/from", getOtherBooks)

	// Start the HTTP server.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Failed to start the server: %s\n", err)
	}
}

// Book example.
type Book struct {
	Title string `json:"title"`
}

func getOtherBooks(w http.ResponseWriter, r *http.Request) {
	//pritn header names 'traceparent'
	println("getOtherBooks::checking traceparent")
	//if r.Header != nil && r.Header["traceparent"] != nil {
	//	println("getOtherBooks::traceparent: " + r.Header["traceparent"][0])
	//} else {
	//	println("getOtherBooks::traceparent: not found")
	//}
	for name, headers := range r.Header {
		for _, value := range headers {
			fmt.Printf("getOtherBooks::%s: %s\n", name, value)
		}
	}

	println("getOtherBooks")
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "Target URL not specified", http.StatusBadRequest)
		return
	}
	println("target = " + target)
	url := "http://" + target + "/books" // Replace with the actual URL.
	println("url = " + url)

	// Make the HTTP GET request.
	//response, err := http.Get(url)
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ctx.WriteString(fmt.Sprintf("Failed to make the HTTP GET request: %s", err.Error()))
		return
	}
	defer response.Body.Close()

	// Check if the response status code is OK (200).
	if response.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadGateway)
		//ctx.WriteString(fmt.Sprintf("HTTP GET request failed with status code: %d", response.StatusCode))
		return
	}

	// Read the response body.
	body, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//ctx.WriteString(fmt.Sprintf("Failed to read response body: %s", err.Error()))
		return
	}

	// Unmarshal the JSON payload into []Books.
	var books []Book
	err = json.Unmarshal(body, &books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//w.WriteString(fmt.Sprintf("Failed to unmarshal JSON payload: %s", err.Error()))
		return
	}
	books = append(books, Book{"Go Go Go"})
	// Send the response body to the client.
	//resp, err := json.Marshal(books)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
	//ctx.Write(resp)
}

func list(w http.ResponseWriter, r *http.Request) {
	println("list::checking traceparent")
	//if r.Header != nil && r.Header["traceparent"] != nil {
	//	println("list::traceparent: " + r.Header["traceparent"][0])
	//} else {
	//	println("list::traceparent: not found")
	//}

	for name, headers := range r.Header {
		for _, value := range headers {
			fmt.Printf("list::%s: %s\n", name, value)
		}
	}

	println("list")
	books := []Book{
		{"Mastering Concurrency in Go"},
		{"Go Design Patterns"},
		{"Black Hat Go"},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}
