package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	// Set up a sample HTTP request
	// Don't use localhost, otherwise, proxy is not working
	req, err := http.NewRequest("GET", "http://server.local:9090", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	fmt.Println(os.Getenv("HTTP_PROXY"))
	// Get the proxy configuration from the environment
	proxyURL, err := http.ProxyFromEnvironment(req)
	if err != nil {
		fmt.Println("Error getting proxy configuration:", err)
		return
	}

	// Print the proxy URL if it's configured
	if proxyURL != nil {
		fmt.Println("Proxy URL:", proxyURL.String())
	} else {
		fmt.Println("Can not find proxy url")
	}

	// Create a custom HTTP client using the proxy configuration
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// Make a request using the custom HTTP client
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status code and content length
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Content Length:", resp.ContentLength)
}
