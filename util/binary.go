package util

import (
	"encoding/binary"
)

func Uint16Bytes(x int) []byte {
	p := make([]byte, 2)
	binary.BigEndian.PutUint16(p, uint16(x))
	return p
}

func Uint32Bytes(x int) []byte {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(x))
	return p
}

func Uint64Bytes(x int) []byte {
	p := make([]byte, 8)
	binary.BigEndian.PutUint64(p, uint64(x))
	return p
}

func UvarintBytes(x int) []byte {
	p := make([]byte, 12)
	n := binary.PutUvarint(p, uint64(x))
	return p[:n]
}

func UvarintSize(x int) int {
	p := make([]byte, 12)
	return binary.PutUvarint(p, uint64(x))
}

func Uint16(p []byte) (int, error) {
	if len(p) < 2 {
		return 0, Error("Not enough bytes")
	}
	x := binary.BigEndian.Uint16(p)
	return int(x), nil
}

func Uint32(p []byte) (int, error) {
	if len(p) < 4 {
		return 0, Error("Not enough bytes")
	}
	x := binary.BigEndian.Uint32(p)
	return int(x), nil
}

func Uint64(p []byte) (int, error) {
	if len(p) < 8 {
		return 0, Error("Not enough bytes")
	}
	x := binary.BigEndian.Uint64(p)
	return int(x), nil
}

func MustUint16(p []byte) int {
	return int(binary.BigEndian.Uint16(p))

}

func MustUint32(p []byte) int {
	return int(binary.BigEndian.Uint32(p))
}

func MustUint64(p []byte) int {
	return int(binary.BigEndian.Uint64(p))
}

func MustUvarint(p []byte) int {
	x, _ := binary.Uvarint(p)
	return int(x)
}

func ReadUint16(r Reader) (int, error) {
	p, err := ReadN(r, 2)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint16(p)
	return int(x), nil
}

func ReadUint32(r Reader) (int, error) {
	p, err := ReadN(r, 4)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint32(p)
	return int(x), nil
}

func ReadUint64(r Reader) (int, error) {
	p, err := ReadN(r, 8)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint64(p)
	return int(x), nil
}

func ReadUvarint(r ByteReader) (int, error) {
	x, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func MustReadUint16(r Reader) int {
	p := MustReadN(r, 2)
	x := binary.BigEndian.Uint16(p)
	return int(x)
}

func MustReadUint32(r Reader) int {
	p := MustReadN(r, 4)
	x := binary.BigEndian.Uint32(p)
	return int(x)
}

func MustReadUint64(r Reader) int {
	p := MustReadN(r, 8)
	x := binary.BigEndian.Uint64(p)
	return int(x)
}

func MustReadUvarint(r ByteReader) int {
	x, err := ReadUvarint(r)
	Check(err)
	return int(x)
}

/*
func FromVarOctet(octet []byte) (p []byte) {
	if j := int(octet[0]); j < 128 {
		p = octet[1:]
	} else {
		j -= 128
		size := int(octet[j])
		fmt.Println(size)
		p = octet[j : j+size]
	}
	return
}

func ReadVarOctet(r Reader) []byte {

}

func ToVarOctet(p []byte) []byte {
	buf := new(bytes.Buffer)
	if size := len(p); size < 128 {
		buf.Write([]byte{uint8(size)})
	} else {
		var i, j int
		var p []byte
		for ; size > 1<<uint(i); i++ {
			p = append([]byte{}, p...)
		}
		if j = i / 8; j*8 != i {
			j++
		}
		buf.Write([]byte{uint8(0x80 | uint(j))})
		bz := make([]byte, j)
		bz[j-1] = uint8(size)
		buf.Write(bz)
	}
	buf.Write(p)
	fmt.Println(buf.Bytes())
	return puf.Bytes()
}
*/
