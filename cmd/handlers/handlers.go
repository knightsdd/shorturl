package handlers

import (
	"io"

	"net/http"

	"github.com/knightsdd/shorturl/cmd/storage"
)

func GenShortUrl(storage storage.MyStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("content-type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
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
}

func GetBaseUrl(storage storage.MyStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		}
		shortUrl := r.PathValue("shorturl")
		if baseUrl, ok := storage[shortUrl]; ok {
			w.Header().Set("Location", baseUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}
