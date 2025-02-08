package main

import (
	"fmt"
	"net"
	"time"
)

func pingHost(ip string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return false
	}

	if conn != nil {
		defer conn.Close()
		fmt.Println("Connection successful!")
		return true
	}

	return false
}

func main() {
	
}
