package util

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/jbenet/go-base58"
)

// Base64
func Base64UrlEncode(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func Base64UrlDecode(b64 string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(b64)
}

func B64Std(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// B58
func BytesToB58(bytes []byte) string {
	return base58.Encode(bytes)
}

func BytesFromB58(b58 string) []byte {
	return base58.Decode(b58)
}

// Hex
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func BytesFromHex(hexstr string) []byte {
	bytes, err := hex.DecodeString(hexstr)
	Check(err)
	return bytes
}

// JSON

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MustMarshalJSON(v interface{}) []byte {
	bytes, err := MarshalJSON(v)
	Check(err)
	return bytes
}

func MarshalIndentJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func MustMarshalIndentJSON(v interface{}) []byte {
	bytes, err := MarshalIndentJSON(v)
	Check(err)
	return bytes
}

func UnmarshalJSON(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}

func MustUnmarshalJSON(bytes []byte, v interface{}) {
	err := UnmarshalJSON(bytes, v)
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
