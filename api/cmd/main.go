// api/cmd/main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Meet Mesh API starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
