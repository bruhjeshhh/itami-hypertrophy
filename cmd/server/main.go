package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
	godotenv.Load()

}
