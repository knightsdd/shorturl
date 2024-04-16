package handlers

import (
	"math/rand"

	"github.com/knightsdd/shorturl/cmd/storage"
)

const (
	defaultLenUrl = 8
)

func randStr() string {
	siqBytes := "abcdifghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, defaultLenUrl)
	for i := range b {
		b[i] = siqBytes[rand.Intn(len(siqBytes))]
	}
	url := string(b)
	return url
}

func generateUrl(storage storage.MyStorage) string {
	for {
		rs := randStr()
		if _, ok := storage[rs]; !ok {
			return rs
		}
	}
}
