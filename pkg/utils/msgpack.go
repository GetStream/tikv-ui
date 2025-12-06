package utils

import (
	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
)

// ParseValue attempts to parse a byte slice as msgpack and return a structured value
// If parsing fails, it returns the raw string
func ParseValue(data []byte) (parsed any, raw string) {
	raw = string(data)

	// If the data is valid UTF-8 and doesn't contain msgpack markers, treat it as plain text
	if isPlainText(data) {
		// Try to decode as JSON even if it looks like plain text
		var decoded any
		if err := json.Unmarshal(data, &decoded); err == nil {
			return decoded, raw
		}
		return raw, raw
	}

	// Try to decode as msgpack
	var decoded any
	if err := msgpack.Unmarshal(data, &decoded); err != nil {
		// If msgpack fails, try JSON
		if err := json.Unmarshal(data, &decoded); err != nil {
			// If both fail, return the raw string as the parsed value
			return raw, raw
		}
		return decoded, raw
	}

	// If msgpack decoded to a simple type (int, float, bool) but the raw data looks like text,
	// it's probably a false positive - return the raw string instead
	switch decoded.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		if isPlainText(data) {
			return raw, raw
		}
	}

	return decoded, raw
}

// isPlainText checks if the data appears to be plain text rather than binary/msgpack
func isPlainText(data []byte) bool {
	// Check if it's valid UTF-8
	if !json.Valid([]byte(`"` + string(data) + `"`)) {
		// If not valid as JSON string, check for msgpack markers
		if len(data) > 0 {
			first := data[0]
			// Check for common msgpack markers (maps, arrays, etc.)
			if first >= 0x80 && first <= 0x9f {
				return false
			}
			if first >= 0xde {
				return false
			}
		}
	}

	// If it contains mostly printable ASCII/UTF-8 characters, treat as plain text
	printable := 0
	for _, b := range data {
		if b >= 32 && b <= 126 || b >= 128 {
			printable++
		}
	}

	// If more than 80% is printable, consider it plain text ( 80% is randomly choosen =) )
	return len(data) > 0 && float64(printable)/float64(len(data)) > 0.8
}
