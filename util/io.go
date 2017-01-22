package util

import (
	"io"
	"io/ioutil"
)

type Reader io.Reader
type Writer io.Writer

func Copy(w io.Writer, r io.Reader) {
	_, err := io.Copy(w, r)
	Check(err)
}

func ReadAll(r io.Reader) []byte {
	bytes, err := ioutil.ReadAll(r)
	Check(err)
	return bytes
}