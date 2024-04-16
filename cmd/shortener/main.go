package main

import (
	"net/http"

	"github.com/knightsdd/shorturl/cmd/handlers"

	"github.com/knightsdd/shorturl/cmd/storage"
)

func main() {
	storage := storage.GetStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", handlers.GenShortUrl(storage))
	mux.HandleFunc("GET /{shorturl}/{$}", handlers.GetBaseUrl(storage))

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
