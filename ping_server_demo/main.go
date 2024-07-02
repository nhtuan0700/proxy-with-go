package main

import (
	"log"
	"net"
	"time"
)

func checkServer(host string, port string) bool {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 10*time.Second)
	if err != nil {
		log.Println("Server is unavailable:", err)
		return false
	}
	defer conn.Close()
	log.Println("Server is reachable")
	return true
}

func main() {
	host := "localhost"
	port := "38081"
	if checkServer(host, port) {
		log.Println("Server is available!")
	} else {
		log.Println("Server is not available.")
	}
}
