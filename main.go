package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "user"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}


func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}