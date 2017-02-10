package spec

import (
	"github.com/kazarena/json-gold/ld"
	ma "github.com/multiformats/go-multiaddr"
	mh "github.com/multiformats/go-multihash"
	"github.com/whyrusleeping/cbor/go"
	. "github.com/zbo14/envoke/common"
)

const (
	IPLD = "ipld"
	JSON = "json_ld"

	LINK_SYMBOL        = "/"
	LINK_TAG    uint64 = 258
)

type Data map[string]interface{}

func CopyData(d Data) Data {
	copy := make(Data)
	for k, v := range d {
		copy[k] = v
	}
	return copy
}

// JSON-LD

func MapJSON(p []byte) Data {
	data := make(Data)
	UnmarshalJSON(p, &data)
	return data
}

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

// Recursive
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
					WrappedObject: maddr.String(), //Bytes()??
				}
			}
		}
	}
	return data
}

type ref struct {
	key    string
	parent int
}

func KeyTrail(i int, refs []ref) []string {
	var keys []string
	for i > 0 {
		ref := refs[i]
		keys = append([]string{ref.key}, keys...)
		i = ref.parent
	}
	return keys
}

func SetInnerValue(data Data, i int, refs []ref, value interface{}) {
	keys := KeyTrail(i, refs)
	var inner interface{} = data
	for i, k := range keys {
		if i == len(keys)-1 {
			inner.(Data)[k] = value
		} else {
			inner = inner.(Data)[k]
		}
	}
}

func GetInnerValue(data Data, i int, refs []ref) (v interface{}) {
	keys := KeyTrail(i, refs)
	var inner interface{} = data
	for i, k := range keys {
		if i == len(keys)-1 {
			v = inner.(Data)[k]
		} else {
			inner = inner.(Data)[k]
		}
	}
	return
}

// Iterative
func IterTransformJSON(data Data) {
	i := 0
	refs := []ref{ref{"", -1}}
	var inner interface{} = data
	for {
		for k, v := range inner.(Data) {
			if _, ok := v.(Data); ok {
				refs = append(refs, ref{k, i})
				continue
			}
			if len(inner.(Data)) == 1 && k == LINK_SYMBOL {
				maddr, err := ma.NewMultiaddr(k + v.(string))
				Check(err)
				value := &cbor.CBORTag{
					Tag:           LINK_TAG,
					WrappedObject: maddr.String(),
				}
				SetInnerValue(data, i, refs, value)
			}
		}
		if i++; i == len(refs) {
			return
		}
		inner = GetInnerValue(data, i, refs)
	}
}

// IPLD

func LinkIPLD(link interface{}) interface{} {
	if link == nil {
		return nil
	}
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
	for k, v := range data {
		if _, ok := v.(Data); ok {
			data[k] = TransformIPLD(v.(Data))
			continue
		}
		if value, ok := v.(*cbor.CBORTag); ok {
			if value.Tag == LINK_TAG {
				str := value.WrappedObject.(string)
				data[k] = Data{LINK_SYMBOL: str[1:]}
			}
		}
	}
	return data
}

func IterTransformIPLD(data Data) {
	i := 0
	refs := []ref{ref{"", -1}}
	var inner interface{} = data
	for {
		for k, v := range inner.(Data) {
			if _, ok := v.(Data); ok {
				refs = append(refs, ref{k, i})
				continue
			}
			if value, ok := v.(*cbor.CBORTag); ok {
				if value.Tag == LINK_TAG {
					str := value.WrappedObject.(string)
					SetInnerValue(data, i, refs, Data{LINK_SYMBOL: str[1:]})
				}
			}
		}
		if i++; i == len(refs) {
			return
		}
		inner = GetInnerValue(data, i, refs)
	}
}

func Multihash(buf []byte, name string) string {
	if _, ok := mh.Names[name]; !ok {
		panic("Unexpected hash func: " + name)
	}
	hash, err := mh.EncodeName(buf, name)
	Check(err)
	return BytesToB58(hash)
}
