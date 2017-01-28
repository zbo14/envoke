package util

import (
	"io"
	"io/ioutil"
)

var EOF = io.EOF

type Reader io.Reader
type Writer io.Writer

type ByteReader io.ByteReader

type MultiReader interface {
	Read([]byte) (int, error)
	ReadByte() (byte, error)
}

func Copy(w io.Writer, r io.Reader) error {
	_, err := io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func MustCopy(w io.Writer, r io.Reader) {
	err := Copy(w, r)
	Check(err)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}

func MustReadAll(r io.Reader) []byte {
	bytes, err := ReadAll(r)
	Check(err)
	return bytes
}

func ReadFull(r io.Reader, buf []byte) error {
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	if size := len(buf); size != n {
		return Errorf("Read %d bytes instead of %d\n", size, n)
	}
	return nil
}

func MustReadFull(r io.Reader, buf []byte) {
	err := ReadFull(r, buf)
	Check(err)
}

func ReadN(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	read, err := io.ReadAtLeast(r, buf, n)
	if err != nil {
		return nil, err
	}
	if read != n {
		return nil, Errorf("Read %d instead of %d bytes", read, n)
	}
	return buf, nil
}

func MustReadN(r io.Reader, n int) []byte {
	bytes, err := ReadN(r, n)
	Check(err)
	return bytes
}

func Write(p []byte, w io.Writer) error {
	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if size := len(p); size != n {
		return Error("Could not write entire slice")
	}
	return nil
}

func MustWrite(p []byte, w io.Writer) {
	err := Write(p, w)
	Check(err)
}
