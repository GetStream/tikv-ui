package utils

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/vmihailenco/msgpack/v5"
)

// ParseValue attempts to parse a byte slice as msgpack / JSON and return a structured value.
func ParseValue(data []byte) (parsed any, raw string) {
	raw = string(data)

	if isPlainText(data) {
		var decoded any
		if err := json.Unmarshal(data, &decoded); err == nil {
			return decoded, raw
		}

		return raw, raw
	}

	var decoded any
	if err := msgpack.Unmarshal(data, &decoded); err != nil {
		if err := json.Unmarshal(data, &decoded); err != nil {
			return raw, raw
		}
		return decoded, raw
	}

	switch decoded.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		if isPlainText(data) {
			return raw, raw
		}
	}

	return decoded, raw
}

// isPlainText checks if the data appears to be plain text (UTF-8, mostly printable)
// rather than binary/msgpack.
func isPlainText(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	if !utf8.Valid(data) {
		return false
	}

	controlCount := 0
	for _, b := range data {
		if b < 0x20 && b != '\n' && b != '\r' && b != '\t' {
			controlCount++
		}
	}

	if controlCount > 0 {
		return false
	}

	printable := 0
	for _, b := range data {
		if (b >= 0x20 && b <= 0x7E) || b >= 0x80 {
			printable++
		}
	}

	return float64(printable)/float64(len(data)) > 0.8
}
