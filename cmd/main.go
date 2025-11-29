package main

import (
	"2025/internal/server"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Test")

	addr := "localhost:8080"
	s := server.NewServer(addr)

	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}
