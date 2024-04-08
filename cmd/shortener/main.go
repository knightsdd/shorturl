package main

import (
	"io"
	"math/rand"
	"net/http"
)

const (
	defaultLenUrl = 8
)

type MyStorage map[string]string

var storage MyStorage = make(MyStorage, 10)

func randStr() string {
	siqBytes := "abcdifghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, defaultLenUrl)
	for i := range b {
		b[i] = siqBytes[rand.Intn(len(siqBytes))]
	}
	url := string(b)
	return url
}

func generateUrl(storage MyStorage) string {
	for {
		rs := randStr()
		if existUrl := storage[rs]; existUrl == "" {
			return rs
		}
	}
}

func getNewUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadGateway)
		return
	}

	shortUrl := generateUrl(storage)
	storage[shortUrl] = string(body)
	schema := "http"
	if r.TLS != nil {
		schema = "https"
	}
	respBody := schema + "://" + r.Host + "/" + shortUrl

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(respBody))
}

func getBaseUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("shorturl")
	if baseUrl := storage[shortUrl]; baseUrl != "" {
		w.Header().Set("Location", baseUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	// storage = make(MyStorage, 10)
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", getNewUrl)
	mux.HandleFunc("GET /{shorturl}/{$}", getBaseUrl)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}
