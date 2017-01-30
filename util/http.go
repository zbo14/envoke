package util

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
)

func HttpGet(url string) (*http.Response, error) {
	return http.Get(url)
}

func HttpPost(url, bodyType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, bodyType, body)
}

func HttpRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	Println(req.Header)
	cli := new(http.Client)
	return cli.Do(req)
}

func HttpsTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: x509.NewCertPool(),
		},
	}
}

func HttpsClient() *http.Client {
	return &http.Client{
		Transport: HttpsTransport(),
	}
}
