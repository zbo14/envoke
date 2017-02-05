package common

import (
	"net/http"
	"net/url"
)

func ParseUrl(rawurl string) (*url.URL, error) {
	return url.Parse(rawurl)
}

func MustParseUrl(rawurl string) *url.URL {
	u, err := ParseUrl(rawurl)
	Check(err)
	return u
}

func ParseQuery(query string) (url.Values, error) {
	return url.ParseQuery(query)
}

func MustParseQuery(query string) url.Values {
	values, err := ParseQuery(query)
	Check(err)
	return values
}

func UrlValues(req *http.Request) (url.Values, error) {
	data, err := ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	vals, err := ParseQuery(string(data))
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func MustUrlValues(req *http.Request) url.Values {
	values, err := UrlValues(req)
	Check(err)
	return values
}
