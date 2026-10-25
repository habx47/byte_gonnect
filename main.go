package main

import (
	"fmt"

	"github.com/habx47/byte_gonnect/listener"
)

func main() {
	port := "8080"
	fmt.Println("Welcome to byte_gonnect")
	fmt.Printf("Listening on port: %s\n", port)
	listener.Listener(port)
}
