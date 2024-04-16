package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/knightsdd/shorturl/cmd/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenShortUrl(t *testing.T) {
	type args struct {
		storage storage.MyStorage
	}
	type want struct {
		contentType string
		code        int
	}
	storage := storage.MyStorage{
		"abcdEFGH": "https://somesite.one",
		"ijklMNOP": "https://somesite.two",
	}
	tests := []struct {
		name       string
		args       args
		requestUrl string
		longUrl    string
		method     string
		want       want
	}{
		{
			name: "positive test to generate short url",
			args: args{
				storage: storage,
			},
			requestUrl: "/",
			longUrl:    "https://testsite.one",
			method:     http.MethodPost,
			want: want{
				contentType: "text/plain",
				code:        201,
			},
		},
		{
			name: "negative test to generate short url",
			args: args{
				storage: storage,
			},
			requestUrl: "/",
			longUrl:    "https://testsite.one",
			method:     http.MethodGet,
			want: want{
				contentType: "text/plain",
				code:        405,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := strings.NewReader(test.longUrl)
			request := httptest.NewRequest(test.method, test.requestUrl, body)
			w := httptest.NewRecorder()

			GenShortUrl(test.args.storage)(w, request)
			result := w.Result()

			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
			assert.Equal(t, test.want.code, result.StatusCode)
			if result.StatusCode == 201 {
				defer result.Body.Close()
				body, err := io.ReadAll(result.Body)

				require.NoError(t, err, "Ошибка при получении тела запроса")

				url := string(body)
				shorturl := url[len(url)-8:]
				testUrl, ok := storage[shorturl]
				require.True(t, ok, "В хранилище нет ссылки")
				assert.Equal(t, test.longUrl, testUrl, "В хранилище неверная ссылка")
			}

		})
	}
}

func TestGetBaseUrl(t *testing.T) {
	type args struct {
		storage storage.MyStorage
	}
	type want struct {
		contentType string
		code        int
		val         string
	}
	shortUrl1 := "abcdEFGH"
	longUrl1 := "https://somesite.one"
	shortUrl2 := "qwerASDF"
	longUrl2 := "https://somesite.two"
	storage := storage.MyStorage{
		"abcdEFGH": "https://somesite.one",
		"qwerASDF": "https://somesite.two",
	}
	tests := []struct {
		name       string
		args       args
		requestUrl string
		method     string
		want       want
	}{
		{
			name: "positive test to get long url 1",
			args: args{
				storage: storage,
			},
			requestUrl: shortUrl1,
			method:     http.MethodGet,
			want: want{
				code: 307,
				val:  longUrl1,
			},
		},
		{
			name: "positive test to get long url 2",
			args: args{
				storage: storage,
			},
			requestUrl: shortUrl2,
			method:     http.MethodGet,
			want: want{
				code: 307,
				val:  longUrl2,
			},
		},
		{
			name: "negative test to get long url 1",
			args: args{
				storage: storage,
			},
			requestUrl: shortUrl2,
			method:     http.MethodPost,
			want: want{
				code: 405,
				val:  "",
			},
		},
		{
			name: "negative test to get long url 2",
			args: args{
				storage: storage,
			},
			requestUrl: "poiuLKJH",
			method:     http.MethodGet,
			want: want{
				code: 400,
				val:  "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/"+test.requestUrl, nil)
			request.SetPathValue("shorturl", test.requestUrl)
			w := httptest.NewRecorder()

			GetBaseUrl(test.args.storage)(w, request)
			result := w.Result()

			assert.Equal(t, test.want.code, result.StatusCode)
			if result.StatusCode == 307 {
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
				require.NotEmpty(t, result.Header.Get("Location"), "Нет заголовка Location")
				assert.Equal(t, test.want.val, result.Header.Get("Location"), "Неверный результат")
			}
		})
	}
}
