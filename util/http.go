package util

import (
	"io"
	"net/http"
)

func HttpClient() *http.Client {
	return new(http.Client)
}

func HttpGet(url string) (*http.Response, error) {
	return http.Get(url)
}

func HttpPost(url, bodyType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, bodyType, body)
}
