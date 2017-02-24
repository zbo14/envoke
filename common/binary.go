package common

import (
	"bytes"
	"encoding/binary"
	"github.com/whyrusleeping/cbor/go"
	"io"
)

const SEG = 1024

// Float32

func WriteFloat32(w io.Writer, x float32) (err error) {
	return binary.Write(w, binary.BigEndian, &x)
}

func WriteFloat32s(w io.Writer, x []float32) (err error) {
	for _, n := range x {
		if err = WriteFloat32(w, n); err != nil {
			return err
		}
	}
	return nil
}

func ReadFloat32(r io.Reader) (x float32, err error) {
	if err = binary.Read(r, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadFloat32s(r io.Reader, seg int) (x []float32, err error) {
	x = make([]float32, seg)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]float32, seg)...)
		}
		x[i], err = ReadFloat32(r)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
	}
}

func ReadNFloat32s(r io.Reader, n int) ([]float32, error) {
	x := make([]float32, n)
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		return nil, err
	}
	return x, nil
}

func BytesFloat32(x float32) []byte {
	buf := new(bytes.Buffer)
	WriteFloat32(buf, x)
	return buf.Bytes()
}

func Float32(p []byte) (float32, error) {
	buf := bytes.NewBuffer(p)
	return ReadFloat32(buf)
}

func BytesFloat32s(x []float32) []byte {
	buf := new(bytes.Buffer)
	WriteFloat32s(buf, x)
	return buf.Bytes()
}

func Float32s(p []byte) ([]float32, error) {
	buf := bytes.NewBuffer(p)
	return ReadFloat32s(buf, SEG)
}

// Float64

func WriteFloat64(w io.Writer, x float64) (err error) {
	return binary.Write(w, binary.BigEndian, &x)
}

func WriteFloat64s(w io.Writer, x []float64) (err error) {
	for _, n := range x {
		if err = WriteFloat64(w, n); err != nil {
			return err
		}
	}
	return nil
}

func ReadFloat64(r io.Reader) (x float64, err error) {
	if err = binary.Read(r, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadFloat64s(r io.Reader, seg int) (x []float64, err error) {
	x = make([]float64, seg)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]float64, seg)...)
		}
		x[i], err = ReadFloat64(r)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
	}
}

func ReadNFloat64s(r io.Reader, n int) ([]float64, error) {
	x := make([]float64, n)
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		return nil, err
	}
	return x, nil
}

func BytesFloat64(x float64) []byte {
	buf := new(bytes.Buffer)
	WriteFloat64(buf, x)
	return buf.Bytes()
}

func Float64(p []byte) (float64, error) {
	buf := bytes.NewBuffer(p)
	return ReadFloat64(buf)
}

func BytesFloat64s(x []float64) []byte {
	buf := new(bytes.Buffer)
	WriteFloat64s(buf, x)
	return buf.Bytes()
}

func Float64s(p []byte) ([]float64, error) {
	buf := bytes.NewBuffer(p)
	return ReadFloat64s(buf, SEG)
}

// Int16

func ReadInt16(r io.Reader) (x int16, err error) {
	if err = binary.Read(r, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadInt16s(r io.Reader, seg int) (x []int16, err error) {
	x = make([]int16, seg)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]int16, seg)...)
		}
		x[i], err = ReadInt16(r)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
	}
}

func ReadNInt16s(r io.Reader, n int) ([]int16, error) {
	x := make([]int16, n)
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		return nil, err
	}
	return x, nil
}

func WriteInt16(w io.Writer, x int16) (err error) {
	return binary.Write(w, binary.BigEndian, &x)
}

func WriteInt16s(w io.Writer, x []int16) (err error) {
	for _, n := range x {
		if err = WriteInt16(w, n); err != nil {
			return err
		}
	}
	return nil
}

func BytesInt16(x int16) []byte {
	buf := new(bytes.Buffer)
	WriteInt16(buf, x)
	p := buf.Bytes()
	return p[len(p)-2:]
}

func Int16(p []byte) (int16, error) {
	buf := bytes.NewBuffer(p)
	return ReadInt16(buf)
}

func BytesInt16s(x []int16) []byte {
	buf := new(bytes.Buffer)
	for _, n := range x {
		buf.Write(BytesInt16(n))
	}
	return buf.Bytes()
}

func Int16s(p []byte) ([]int16, error) {
	buf := bytes.NewBuffer(p)
	return ReadInt16s(buf, SEG)
}

// Int32

func ReadInt32(r io.Reader) (x int32, err error) {
	if err = binary.Read(r, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadInt32s(r io.Reader, seg int) (x []int32, err error) {
	x = make([]int32, seg)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]int32, seg)...)
		}
		x[i], err = ReadInt32(r)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
	}
}

func ReadNInt32s(r io.Reader, n int) ([]int32, error) {
	x := make([]int32, n)
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		return nil, err
	}
	return x, nil
}

func WriteInt32(w io.Writer, x int32) (err error) {
	return binary.Write(w, binary.BigEndian, &x)
}

func WriteInt32s(w io.Writer, x []int32) (err error) {
	for _, n := range x {
		if err = WriteInt32(w, n); err != nil {
			return err
		}
	}
	return nil
}

func BytesInt32(x int32) []byte {
	buf := new(bytes.Buffer)
	WriteInt32(buf, x)
	p := buf.Bytes()
	return p[len(p)-4:]
}

func Int32(p []byte) (x int32, err error) {
	buf := bytes.NewBuffer(p)
	return ReadInt32(buf)
}

func BytesInt32s(x []int32) []byte {
	buf := new(bytes.Buffer)
	for _, n := range x {
		buf.Write(BytesInt32(n))
	}
	return buf.Bytes()
}

func Int32s(p []byte) ([]int32, error) {
	buf := bytes.NewBuffer(p)
	return ReadInt32s(buf, SEG)
}

// Int64

func ReadInt64(r io.Reader) (x int64, err error) {
	if err = binary.Read(r, binary.BigEndian, &x); err != nil {
		return 0, err
	}
	return x, nil
}

func ReadInt64s(r io.Reader, seg int) (x []int64, err error) {
	x = make([]int64, seg)
	for i := 0; ; i++ {
		if i == len(x) {
			x = append(x, make([]int64, seg)...)
		}
		x[i], err = ReadInt64(r)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return x[:i], nil
			}
			return nil, err
		}
	}
}

func ReadNInt64s(r io.Reader, n int) ([]int64, error) {
	x := make([]int64, n)
	if err := binary.Read(r, binary.BigEndian, &x); err != nil {
		return nil, err
	}
	return x, nil
}

func WriteInt64(w io.Writer, x int64) (err error) {
	return binary.Write(w, binary.BigEndian, &x)
}

func WriteInt64s(w io.Writer, x []int64) (err error) {
	for _, n := range x {
		if err = WriteInt64(w, n); err != nil {
			return err
		}
	}
	return nil
}

func BytesInt64(x int64) []byte {
	buf := new(bytes.Buffer)
	WriteInt64(buf, x)
	p := buf.Bytes()
	return p[len(p)-8:]
}

func Int64(p []byte) (int64, error) {
	buf := bytes.NewBuffer(p)
	return ReadInt64(buf)
}

func BytesInt64s(x []int64) []byte {
	buf := new(bytes.Buffer)
	for _, n := range x {
		buf.Write(BytesInt64(n))
	}
	return buf.Bytes()
}

func Int64s(p []byte) ([]int64, error) {
	buf := bytes.NewBuffer(p)
	return ReadInt64s(buf, SEG)
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
