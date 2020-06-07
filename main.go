package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port += "8080"
	}
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World!")
	}))
}
