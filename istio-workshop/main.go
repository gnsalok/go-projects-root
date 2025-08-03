package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	version := os.Getenv("VERSION")
	if version == "" {
		version = "v1"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go app! Version: %s\n", version)
	})
	http.ListenAndServe(":8080", nil)
}
