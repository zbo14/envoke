package common

import (
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

const MAX_MEMORY int64 = 1000000000000

func MultipartForm(req *http.Request) (*multipart.Form, error) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, Error("Expected mimetype=multipart; got mimetype=" + mediaType)
	}
	r := multipart.NewReader(req.Body, params["boundary"])
	return r.ReadForm(MAX_MEMORY)
}

func MustMultipartForm(req *http.Request) *multipart.Form {
	form, err := MultipartForm(req)
	Check(err)
	return form
}
