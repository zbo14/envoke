package rrm

import (
	"encoding/json"
	"github.com/kazarena/json-gold/ld"
)

// REGEX
const LAT_LONG = `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?),\s*[-+]?(180(\.0+)?|((1[0-7]\d)|([1-9]?\d))(\.\d+)?)$`

// json to map[string]interface{} for json-ld interpretation
func MapJSON(data []byte) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	err := json.Unmarshal(data, &mp)
	if err != nil {
		return nil, err
	}
	return mp, nil
}

// json-ld methods
func CompactJSON(data []byte) (map[string]interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	mp, err := MapJSON(data)
	if err != nil {
		return nil, err
	}
	output, err := proc.Compact(mp, nil, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func ExpandJSON(data []byte) ([]interface{}, error) {
	proc := ld.NewJsonLdProcessor()
	mp, err := MapJSON(data)
	if err != nil {
		return nil, err
	}
	output, err := proc.Expand(mp, nil)
	if err != nil {
		return nil, err
	}
	return output, nil
}
