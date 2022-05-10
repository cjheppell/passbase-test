package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	register(mux)
	http.ListenAndServe(":8081", mux)
}
