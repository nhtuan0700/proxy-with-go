package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func sendRequest(method string, url string, data string) {
	// Set up a sample HTTP request
	// Don't use localhost, otherwise, proxy is not working
	req, err := http.NewRequest(method, url, strings.NewReader(data))
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
		// This is default transport
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

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read body")
	}
	// Print the response status code and content length
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Content Length:", resp.ContentLength)
	fmt.Println("Response data:", string(d))
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	sendRequest("GET", "http://server.local:28080", "Tuan")
}
