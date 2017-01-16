package util

import "github.com/tendermint/go-wire"

const (
	QUERY_KEY   byte = 1
	QUERY_INDEX byte = 2
	QUERY_SIZE  byte = 3
)

func EmptyQuery(QueryType byte) []byte {
	return []byte{QueryType}
}

func KeyQuery(key []byte) []byte {
	query := make([]byte, wire.ByteSliceSize(key)+1)
	buf := query
	buf[0] = QUERY_KEY
	buf = buf[1:]
	wire.PutByteSlice(buf, key)
	return query
}

func IndexQuery(i int) []byte {
	query := make([]byte, 100)
	buf := query
	buf[0] = QUERY_INDEX
	buf = buf[1:]
	n, err := wire.PutVarint(buf, i)
	if err != nil {
		return nil
	}
	query = query[:n+1]
	return query
}
