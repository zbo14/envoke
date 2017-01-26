package util

import (
	"encoding/binary"
)

func Uint16Bytes(x uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, x)
	return b
}

func Uint32Bytes(x uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, x)
	return b
}

func Uint64Bytes(x uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, x)
	return b
}

func UvarintBytes(x uint64) []byte {
	b := make([]byte, 8)
	binary.PutUvarint(b, x)
	return b
}

func Uint16(b []byte) (uint16, error) {
	if len(b) < 2 {
		return 0, Error("Not enough bytes")
	}
	return binary.BigEndian.Uint16(b), nil
}

func Uint32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, Error("Not enough bytes")
	}
	return binary.BigEndian.Uint32(b), nil
}

func Uint64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, Error("Not enough bytes")
	}
	return binary.BigEndian.Uint64(b), nil
}

func Uvarint(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, Error("Not enough bytes")
	}
	x, _ := binary.Uvarint(b)
	return x, nil
}

func MustUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func MustUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func MustUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func ReadUint16(r Reader) (uint16, error) {
	b, err := ReadN(r, 2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b), nil
}

func ReadUint32(r Reader) (uint32, error) {
	b, err := ReadN(r, 4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

func ReadUint64(r Reader) (uint64, error) {
	b, err := ReadN(r, 8)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

func MustReadUint16(r Reader) uint16 {
	b := MustReadN(r, 2)
	return binary.BigEndian.Uint16(b)
}

func MustReadUint32(r Reader) uint32 {
	b := MustReadN(r, 2)
	return binary.BigEndian.Uint32(b)
}

func MustReadUint64(r Reader) uint64 {
	b := MustReadN(r, 2)
	return binary.BigEndian.Uint64(b)
}
