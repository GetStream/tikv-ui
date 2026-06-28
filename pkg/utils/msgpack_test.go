package utils

import (
	"encoding/base64"
	"testing"

	"github.com/vmihailenco/msgpack/v5"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantRaw  string
		wantJSON bool // if true, expect parsed to be a structured object/slice
	}{
		{
			name:     "Simple JSON Object",
			input:    []byte(`{"foo":"bar"}`),
			wantRaw:  `{"foo":"bar"}`,
			wantJSON: true,
		},
		{
			name:     "JSON Array",
			input:    []byte(`[1, 2, 3]`),
			wantRaw:  `[1, 2, 3]`,
			wantJSON: true,
		},
		{
			name:     "Plain String",
			input:    []byte(`Hello World`),
			wantRaw:  `Hello World`,
			wantJSON: false,
		},
		{
			name:     "JSON String",
			input:    []byte(`"quoted string"`),
			wantRaw:  `"quoted string"`,
			wantJSON: true,
		},
		{
			name: "Spacing",
			input: []byte(`{
"foo":"bar"
}`),
			wantRaw: `{
"foo":"bar"
}`,
			wantJSON: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, raw := ParseValue(tt.input)
			if raw != tt.wantRaw {
				t.Errorf("ParseValue() raw = %v, want %v", raw, tt.wantRaw)
			}

			if tt.wantJSON {
				switch parsed.(type) {
				case map[string]any, []any, string, float64, bool, nil:
				default:
					t.Errorf("ParseValue() parsed type = %T, want JSON type", parsed)
				}

				if tt.name == "Simple JSON Object" {
					m, ok := parsed.(map[string]any)
					if !ok || m["foo"] != "bar" {
						t.Errorf("Failed to parse JSON object correctly: %v", parsed)
					}
				}
			} else {
				if parsed != raw {
					t.Errorf("ParseValue() parsed = %v, want %v", parsed, raw)
				}
			}
		})
	}
}

func TestFormatRawValue_Base64(t *testing.T) {
	data := []byte{0x00, 0x01, 0xff}
	want := base64.StdEncoding.EncodeToString(data)
	if got := FormatRawValue(data); got != want {
		t.Fatalf("FormatRawValue() = %q, want %q", got, want)
	}
}

func makeFeedActivity() map[string]any {
	return map[string]any{
		"fid": "user:test",
		"op":  int8(1),
		"a":   "$2b454049-e5bd-4efb-9c55-395beb3ad026",
		"u":   "test",
		"o":   "user:test",
		"v":   int8(1),
		"c":   &msgpackExt{Payload: []byte{0x18, 0, 0, 0, 0x59, 0, 0, 0}},
		"ag":  []string{"fg:notification:0", "_post_2026-05-22", "fg:stories:0", "test"},
		"at":  "post",
	}
}

func TestParseValue_FeedsMsgpackActivity(t *testing.T) {
	ensureMsgpackExt()

	data, err := msgpack.Marshal(makeFeedActivity())
	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := ParseValue(data)
	m, ok := parsed.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", parsed)
	}
	if m["fid"] != "user:test" {
		t.Fatalf("fid = %v", m["fid"])
	}
	if m["at"] != "post" {
		t.Fatalf("at = %v", m["at"])
	}

	switch m["c"].(type) {
	case int64, int, uint64, map[string]any:
	default:
		t.Fatalf("unexpected c type %T: %v", m["c"], m["c"])
	}
}

func TestParseValue_FeedsMsgpackActivityArray(t *testing.T) {
	ensureMsgpackExt()

	activities := []any{makeFeedActivity(), makeFeedActivity()}
	data, err := msgpack.Marshal(activities)
	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := ParseValue(data)
	arr, ok := parsed.([]any)
	if !ok {
		t.Fatalf("expected array, got %T", parsed)
	}
	if len(arr) != 2 {
		t.Fatalf("len = %d", len(arr))
	}
}

func TestParseValue_FeedsNestedPayload(t *testing.T) {
	ensureMsgpackExt()

	ts := []byte{0x18, 0, 0, 0, 0x59, 0x4d, 0x60, 0x00}
	inner, err := msgpack.Marshal(map[string]any{
		"a": "test",
		"c": &msgpackExt{Payload: ts},
		"g": "user",
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := msgpack.Marshal(map[string]any{
		"c": &msgpackExt{Payload: ts},
		"p": inner,
		"a": "test",
		"g": "user",
		"v": &msgpackExt{Payload: []byte{0xd0, 0x90, 0xd4, 0x87, 0x0c}},
		"i": []byte{0x01, 0x02, 0x03},
	})
	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := ParseValue(data)
	m, ok := parsed.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", parsed)
	}

	p, ok := m["p"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map in p, got %T: %v", m["p"], m["p"])
	}
	if p["a"] != "test" || p["g"] != "user" {
		t.Fatalf("nested p = %v", p)
	}

	i, ok := m["i"].(map[string]any)
	if !ok {
		t.Fatalf("expected hex wrapper for opaque i, got %T", m["i"])
	}
	if i["$hex"] != "010203" {
		t.Fatalf("i hex = %v", i["$hex"])
	}
}

func TestParseValue_ConcatenatedMsgpackMaps(t *testing.T) {
	ensureMsgpackExt()

	a1, err := msgpack.Marshal(makeFeedActivity())
	if err != nil {
		t.Fatal(err)
	}
	a2, err := msgpack.Marshal(makeFeedActivity())
	if err != nil {
		t.Fatal(err)
	}
	data := append(a1, a2...)

	parsed, _ := ParseValue(data)
	arr, ok := parsed.([]any)
	if !ok {
		t.Fatalf("expected array of maps, got %T", parsed)
	}
	if len(arr) != 2 {
		t.Fatalf("len = %d", len(arr))
	}
}

func TestParseValue_VersionedWrapper(t *testing.T) {
	ensureMsgpackExt()

	inner, err := msgpack.Marshal(map[string]any{
		"fid": "user:test",
		"op":  int8(1),
		"at":  "post",
	})
	if err != nil {
		t.Fatal(err)
	}
	// Feeds storage: leading version 0 (fixint) + msgpack payload.
	data := append([]byte{0x00}, inner...)

	parsed, _ := ParseValue(data)
	m, ok := parsed.(map[string]any)
	if !ok {
		t.Fatalf("expected unwrapped map, got %T: %v", parsed, parsed)
	}
	if m["fid"] != "user:test" {
		t.Fatalf("fid = %v", m["fid"])
	}
}
