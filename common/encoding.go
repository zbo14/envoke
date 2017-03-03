package common

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"github.com/jbenet/go-base58"
	"io"
)

// Base64
func Base64UrlEncode(p []byte) string {
	return base64.RawURLEncoding.EncodeToString(p)
}

func Base64UrlDecode(b64 string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(b64)
}

func MustBase64UrlDecode(b64 string) []byte {
	p, err := Base64UrlDecode(b64)
	Check(err)
	return p
}

func Base64StdEncode(p []byte) string {
	return base64.StdEncoding.EncodeToString(p)
}

func Base64StdDecode(b64 string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64)
}

func MustBase64StdDecode(b64 string) []byte {
	p, err := Base64StdDecode(b64)
	Check(err)
	return p
}

// B58
func BytesToB58(p []byte) string {
	return base58.Encode(p)
}

func BytesFromB58(b58 string) []byte {
	return base58.Decode(b58)
}

// Hex
func BytesToHex(p []byte) string {
	return hex.EncodeToString(p)
}

func BytesFromHex(hexstr string) []byte {
	p, err := hex.DecodeString(hexstr)
	Check(err)
	return p
}

// JSON
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MustMarshalJSON(v interface{}) []byte {
	p, err := MarshalJSON(v)
	Check(err)
	return p
}

func MarshalIndentJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func MustMarshalIndentJSON(v interface{}) []byte {
	p, err := MarshalIndentJSON(v)
	Check(err)
	return p
}

func UnmarshalJSON(p []byte, v interface{}) error {
	return json.Unmarshal(p, v)
}

func MustUnmarshalJSON(p []byte, v interface{}) {
	err := UnmarshalJSON(p, v)
	Check(err)
}

func ReadJSON(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

func MustReadJSON(r io.Reader, v interface{}) {
	err := ReadJSON(r, v)
	Check(err)
}

func WriteJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func MustWriteJSON(w io.Writer, v interface{}) {
	err := WriteJSON(w, v)
	Check(err)
}

// PEM
func BlockPEM(p []byte, _type string) *pem.Block {
	return &pem.Block{
		Bytes: p,
		Type:  _type,
	}
}

func EncodePEM(b *pem.Block) []byte {
	return pem.EncodeToMemory(b)
}

func DecodePEM(p []byte) (*pem.Block, []byte) {
	return pem.Decode(p)
}
