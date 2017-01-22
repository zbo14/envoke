package util

import (
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

const MAX_MEMORY int64 = 1000000000000

func UrlValues(req *http.Request) (url.Values, error) {
	data := ReadAll(req.Body)
	vals, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func MultipartForm(req *http.Request) (*multipart.Form, error) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, Error("Expected mimetype=multipart; got mimetype=" + mediaType)
	}
	mr := multipart.NewReader(req.Body, params["boundary"])
	return mr.ReadForm(MAX_MEMORY)
}
