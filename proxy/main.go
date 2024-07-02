package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Print details of incoming request
	fmt.Printf("Received request for: [%s] %s\n", r.Method, r.URL.String())

	// Create a new HTTP client
	client := &http.Client{}

	// Forward the request to the destination server
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}
	if err != nil {
		http.Error(w, "Failed to init request", http.StatusBadRequest)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy headers from the response to the client
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code of the response
	w.WriteHeader(resp.StatusCode)

	// Copy the response body to the client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Create a new HTTP server
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	http.HandleFunc("/", proxyHandler)
	addr := os.Getenv("HTTP_PROXY")

	// Start the server on port 8080
	fmt.Printf("Proxy server listening on %s ...\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
