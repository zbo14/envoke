package common

import (
	"bytes"
	"encoding/binary"
	"github.com/whyrusleeping/cbor/go"
)

// Int16

func Int16Bytes(x int16) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &x)
	p := buf.Bytes()
	return p[len(p)-2:]
}

func Int16(p []byte) (x int16, err error) {
	if len(p) < 2 {
		return 0, ErrInvalidSize
	}
	buf := bytes.NewBuffer(p)
	if err = binary.Read(buf, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func Int16SliceBytes(x []int16) (p []byte) {
	buf := new(bytes.Buffer)
	for _, n := range x {
		binary.Write(buf, binary.BigEndian, &n)
		q := buf.Bytes()
		p = append(p, q[len(q)-2:]...)
	}
	return p
}

func Int16Slice(p []byte) ([]int16, error) {
	if len(p) < 2 || len(p)%2 != 0 {
		return nil, ErrInvalidSize
	}
	buf := bytes.NewBuffer(p)
	x := make([]int16, len(p)/2)
	for i := range x {
		if err := binary.Read(buf, binary.BigEndian, &x[i]); err != nil {
			return nil, err
		}
	}
	return x, nil
}

// Int32

func Int32Bytes(x int32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, &x)
	p := buf.Bytes()
	return p[len(p)-4:]
}

func Int32(p []byte) (x int32, err error) {
	if len(p) < 4 {
		return 0, ErrInvalidSize
	}
	buf := bytes.NewBuffer(p)
	if err = binary.Read(buf, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func Int32SliceBytes(x []int32) (p []byte) {
	buf := new(bytes.Buffer)
	for _, n := range x {
		binary.Write(buf, binary.BigEndian, &n)
		q := buf.Bytes()
		p = append(p, q[len(q)-4:]...)
	}
	return p
}

func Int32Slice(p []byte) ([]int32, error) {
	if len(p) < 4 || len(p)%4 != 0 {
		return nil, ErrInvalidSize
	}
	buf := bytes.NewBuffer(p)
	x := make([]int32, len(p)/4)
	for i := range x {
		if err := binary.Read(buf, binary.BigEndian, &x[i]); err != nil {
			return nil, err
		}
	}
	return x, nil
}

// Uint16

func Uint16Bytes(x int) []byte {
	p := make([]byte, 2)
	binary.BigEndian.PutUint16(p, uint16(x))
	return p
}

func Uint16(p []byte) (int, error) {
	if len(p) < 2 {
		return 0, ErrInvalidSize
	}
	x := binary.BigEndian.Uint16(p)
	return int(x), nil
}

func MustUint16(p []byte) int {
	return int(binary.BigEndian.Uint16(p))

}

func ReadUint16(r Reader) (int, error) {
	p, err := ReadN(r, 2)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint16(p)
	return int(x), nil
}

func MustReadUint16(r Reader) int {
	p := MustReadN(r, 2)
	x := binary.BigEndian.Uint16(p)
	return int(x)
}

// Uint32

func Uint32Bytes(x int) []byte {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(x))
	return p
}

func Uint32(p []byte) (int, error) {
	if len(p) < 4 {
		return 0, ErrInvalidSize
	}
	x := binary.BigEndian.Uint32(p)
	return int(x), nil
}

func MustUint32(p []byte) int {
	return int(binary.BigEndian.Uint32(p))
}

func ReadUint32(r Reader) (int, error) {
	p, err := ReadN(r, 4)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint32(p)
	return int(x), nil
}

func MustReadUint32(r Reader) int {
	p := MustReadN(r, 4)
	x := binary.BigEndian.Uint32(p)
	return int(x)
}

// Uint64

func Uint64Bytes(x int) []byte {
	p := make([]byte, 8)
	binary.BigEndian.PutUint64(p, uint64(x))
	return p
}

func Uint64(p []byte) (int, error) {
	if len(p) < 8 {
		return 0, ErrInvalidSize
	}
	x := binary.BigEndian.Uint64(p)
	return int(x), nil
}

func MustUint64(p []byte) int {
	return int(binary.BigEndian.Uint64(p))
}

func ReadUint64(r Reader) (int, error) {
	p, err := ReadN(r, 8)
	if err != nil {
		return 0, err
	}
	x := binary.BigEndian.Uint64(p)
	return int(x), nil
}

func MustReadUint64(r Reader) int {
	p := MustReadN(r, 8)
	x := binary.BigEndian.Uint64(p)
	return int(x)
}

// Uvarint

func UvarintBytes(x int) []byte {
	p := make([]byte, 12)
	n := binary.PutUvarint(p, uint64(x))
	return p[:n]
}

func UvarintSize(x int) int {
	p := make([]byte, 12)
	return binary.PutUvarint(p, uint64(x))
}

func MustUvarint(p []byte) int {
	x, _ := binary.Uvarint(p)
	return int(x)
}

func ReadUvarint(r ByteReader) (int, error) {
	x, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

func MustReadUvarint(r ByteReader) int {
	x, err := ReadUvarint(r)
	Check(err)
	return int(x)
}

// VarBytes

func VarBytes(p []byte) []byte {
	size := UvarintBytes(len(p))
	return append(size, p...)
}

func ReadVarBytes(r MyReader) ([]byte, error) {
	n, err := ReadUvarint(r)
	if err != nil {
		return nil, err
	}
	return ReadN(r, n)
}

func MustReadVarBytes(r MyReader) []byte {
	p, err := ReadVarBytes(r)
	Check(err)
	return p
}

func WriteVarBytes(p []byte, w Writer) error {
	v := VarBytes(p)
	return Write(v, w)
}

func MustWriteVarBytes(p []byte, w Writer) {
	err := WriteVarBytes(p, w)
	Check(err)
}

/*

// Octet

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

// CBOR

func DumpCBOR(v interface{}) ([]byte, error) {
	return cbor.Dumps(v)
}

func MustDumpCBOR(v interface{}) []byte {
	p, err := DumpCBOR(v)
	Check(err)
	return p
}

func LoadCBOR(p []byte, v interface{}) error {
	return cbor.Loads(p, v)
}

func MustLoadCBOR(p []byte, v interface{}) {
	err := LoadCBOR(p, v)
	Check(err)
}
