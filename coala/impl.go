package coala

import (
	"github.com/kazarena/json-gold/ld"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/whyrusleeping/cbor/go"
	. "github.com/zballs/envoke/util"
)

const (
	LINK_SYMBOL        = "/"
	LINK_TAG    uint64 = 258
)

// JSON-ld

// json to map[string]interface{} for json-ld interpretation
func MapJSON(p []byte) Data {
	data := make(Data)
	UnmarshalJSON(p, &data)
	return data
}

// json-ld methods
func CompactJSON(p []byte) (Data, error) {
	proc := ld.NewJsonLdProcessor()
	data := MapJSON(p)
	output, err := proc.Compact(data, nil, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func MustCompactJSON(p []byte) Data {
	output, err := CompactJSON(p)
	Check(err)
	return output
}

func ExpandJSON(p []byte) ([]interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	data := MapJSON(p)
	output, err := proc.Expand(data, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func MustExpandJSON(p []byte) []interface{} {
	output, err := ExpandJSON(p)
	Check(err)
	return output
}

func TransformJSON(data Data) interface{} {
	for key, value := range data {
		switch value.(type) {
		case Data:
			data[key] = TransformJSON(value.(Data))
		default:
			if len(data) == 1 && key == LINK_SYMBOL {
				maddr, err := ma.NewMultiaddr(key + value.(string))
				Check(err)
				return &cbor.CBORTag{
					Tag:           LINK_TAG,
					WrappedObject: maddr.Bytes(), //String()??
				}
			}
		}
	}
	return data
}

// IPLD

func LinkIPLD(link interface{}) interface{} {
	data := make(Data)
	data[LINK_SYMBOL] = link
	return data
}

func LinksIPLD(links ...interface{}) []interface{} {
	numLinks := len(links)
	if numLinks <= 1 {
		panic("Expected more than one link")
	}
	datas := make([]interface{}, numLinks)
	for i, _ := range datas {
		datas[i] = LinkIPLD(links[i])
	}
	return datas
}

func EncodeIPLD(data Data) []byte {
	return MustDumpCBOR(data)
}

func DecodeIPLD(p []byte) (Data, error) {
	data := make(Data)
	if err := LoadCBOR(p, data); err != nil {
		return nil, err
	}
	return data, nil
}

func TransformIPLD(data Data) interface{} {
	for key, value := range data {
		switch value.(type) {
		case Data:
			data[key] = TransformIPLD(data)
		case *cbor.CBORTag:
			v := value.(*cbor.CBORTag)
			if v.Tag == LINK_TAG {
				p := v.WrappedObject.([]byte)
				return Data{LINK_SYMBOL: p}
			}
		}
	}
	return data
}

func Multihash(buf []byte, name string) string {
	if _, ok := mh.Names[name]; !ok {
		panic("Unexpected hash func: " + name)
	}
	hash, err := mh.EncodeName(buf, name)
	Check(err)
	return BytesToB58(hash)
}
