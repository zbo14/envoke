package util

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jbenet/go-base58"
)

// Base64
func Base64UrlEncode(p []byte) string {
	return base64.RawURLEncoding.EncodeToString(p)
}

func Base64UrlDecode(b64 string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(b64)
}

func B64Std(p []byte) string {
	return base64.StdEncoding.EncodeToString(p)
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

func ReadJSON(r Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

func MustReadJSON(r Reader, v interface{}) {
	err := ReadJSON(r, v)
	Check(err)
}

func WriteJSON(w Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v)
}

func MustWriteJSON(w Writer, v interface{}) {
	err := WriteJSON(w, v)
	Check(err)
}
