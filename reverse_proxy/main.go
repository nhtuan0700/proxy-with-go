package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/nhtuan0700/proxy-with-go/reverse_proxy/load_balancer"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	targets := []string{
		os.Getenv("HOST_1"),
		os.Getenv("HOST_2"),
		os.Getenv("HOST_3"),
	}

	loadBalancer, err := load_balancer.NewLoadBalancer(targets)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", loadBalancer.ServeHTTP)

	host := os.Getenv("REVERSE_PROXY_URL")
	fmt.Println("Reverse proxy is running on", host)
	http.ListenAndServe(host, nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutting down server...")
}
