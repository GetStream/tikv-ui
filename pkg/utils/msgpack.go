package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"sync"
	"unicode/utf8"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// msgpackExt handles Stream feeds extension type 0 (e.g. the compact "c" timestamp field).
type msgpackExt struct {
	Payload []byte
}

func (e *msgpackExt) MarshalMsgpack() ([]byte, error) {
	return e.Payload, nil
}

func (e *msgpackExt) UnmarshalMsgpack(b []byte) error {
	e.Payload = append([]byte(nil), b...)
	return nil
}

var registerMsgpackExt sync.Once

func ensureMsgpackExt() {
	registerMsgpackExt.Do(func() {
		// Feeds TiKV values use msgpack ext types for compact binary fields (e.g. "c", "v").
		for extID := int8(0); extID < 32; extID++ {
			msgpack.RegisterExt(extID, (*msgpackExt)(nil))
		}
	})
}

// FormatRawValue returns a JSON-safe representation of raw TiKV bytes.
func FormatRawValue(data []byte) string {
	if isPlainText(data) {
		return string(data)
	}
	return base64.StdEncoding.EncodeToString(data)
}

// ParseValue attempts to parse a byte slice as msgpack / JSON and return a structured value.
func ParseValue(data []byte) (parsed any, raw string) {
	raw = FormatRawValue(data)

	if isPlainText(data) {
		var decoded any
		if err := json.Unmarshal(data, &decoded); err == nil {
			return decoded, raw
		}

		return raw, raw
	}

	decoded, err := decodeMsgpack(data)
	if err != nil {
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

	return unwrapVersioned(decoded), raw
}

func decodeMsgpack(data []byte) (any, error) {
	ensureMsgpackExt()

	values, err := decodeMsgpackValues(data)
	if err != nil {
		return nil, err
	}

	switch len(values) {
	case 0:
		return nil, io.EOF
	case 1:
		return values[0], nil
	default:
		return values, nil
	}
}

func decodeMsgpackValues(data []byte) ([]any, error) {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	dec.UsePreallocateValues(true)

	var values []any
	for {
		var v any
		err := dec.Decode(&v)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		values = append(values, normalizeDecoded(v))
	}
	return values, nil
}

// unwrapVersioned strips the feeds storage version prefix: a leading 0 followed by the payload.
func unwrapVersioned(v any) any {
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		return v
	}
	if version, ok := arr[0].(int64); ok && version == 0 {
		return arr[1]
	}
	if version, ok := arr[0].(int8); ok && version == 0 {
		return arr[1]
	}
	if version, ok := arr[0].(uint8); ok && version == 0 {
		return arr[1]
	}
	if version, ok := arr[0].(int); ok && version == 0 {
		return arr[1]
	}
	return v
}

func normalizeDecoded(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			out[k] = normalizeDecoded(val)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			out[normalizeMapKey(k)] = normalizeDecoded(val)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = normalizeDecoded(val)
		}
		return out
	case *msgpackExt:
		if x == nil {
			return nil
		}
		return extToJSON(0, x.Payload)
	case msgpackExt:
		return extToJSON(0, x.Payload)
	case []byte:
		if nested, ok := tryDecodeNestedMsgpack(x); ok {
			return nested
		}
		return bytesToJSON(x)
	default:
		return v
	}
}

// tryDecodeNestedMsgpack decodes bin fields that embed msgpack values (e.g. feeds "p").
func tryDecodeNestedMsgpack(data []byte) (any, bool) {
	if !looksLikeEmbeddedMsgpack(data) {
		return nil, false
	}

	values, err := decodeMsgpackValues(data)
	if err != nil || len(values) == 0 {
		return nil, false
	}
	if len(values) == 1 {
		return values[0], true
	}
	return values, true
}

func looksLikeEmbeddedMsgpack(data []byte) bool {
	if len(data) < 2 {
		return false
	}

	c := data[0]
	if msgpcode.IsFixedMap(c) {
		return true
	}
	if msgpcode.IsFixedString(c) {
		return true
	}
	if msgpcode.IsFixedArray(c) && c > msgpcode.FixedArrayLow {
		return true
	}

	switch c {
	case msgpcode.Map16, msgpcode.Map32,
		msgpcode.Array16, msgpcode.Array32,
		msgpcode.Str8, msgpcode.Str16, msgpcode.Str32,
		msgpcode.Bin8, msgpcode.Bin16, msgpcode.Bin32,
		msgpcode.FixExt1, msgpcode.FixExt2, msgpcode.FixExt4, msgpcode.FixExt8, msgpcode.FixExt16,
		msgpcode.Ext8, msgpcode.Ext16, msgpcode.Ext32:
		return true
	}

	return false
}

func bytesToJSON(b []byte) map[string]any {
	return map[string]any{
		"$hex": hex.EncodeToString(b),
		"$len": len(b),
	}
}

func normalizeMapKey(k any) string {
	switch key := k.(type) {
	case string:
		return key
	case []byte:
		return string(key)
	case int:
		return strconv.Itoa(key)
	case int8:
		return strconv.FormatInt(int64(key), 10)
	case int16:
		return strconv.FormatInt(int64(key), 10)
	case int32:
		return strconv.FormatInt(int64(key), 10)
	case int64:
		return strconv.FormatInt(key, 10)
	case uint:
		return strconv.FormatUint(uint64(key), 10)
	case uint8:
		return strconv.FormatUint(uint64(key), 10)
	case uint16:
		return strconv.FormatUint(uint64(key), 10)
	case uint32:
		return strconv.FormatUint(uint64(key), 10)
	case uint64:
		return strconv.FormatUint(key, 10)
	default:
		return fmt.Sprint(key)
	}
}

func extToJSON(extID int8, payload []byte) any {
	if ts, ok := decodeExtTimestamp(payload); ok {
		return ts
	}
	return map[string]any{
		"$ext": extID,
		"$hex": hex.EncodeToString(payload),
		"$len": len(payload),
	}
}

// decodeExtTimestamp interprets 8-byte feeds ext payloads as nanosecond timestamps when plausible.
func decodeExtTimestamp(payload []byte) (int64, bool) {
	if len(payload) != 8 {
		return 0, false
	}
	for _, n := range []uint64{
		binary.BigEndian.Uint64(payload),
		binary.LittleEndian.Uint64(payload),
	} {
		if n >= 1e17 && n <= 9e18 {
			ts := int64(n)
			if ts > 0 {
				return ts, true
			}
		}
	}
	return 0, false
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
