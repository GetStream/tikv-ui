package utils

import (
	"testing"
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
			wantJSON: true, // It is a valid JSON string, so it should be parsed as string
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
				// Assert that parsed is NOT the raw string (unless the JSON value IS a string)
				// But specifically for object/array it should be map/slice
				switch parsed.(type) {
				case map[string]any, []any, string, float64, bool, nil:
					// OK
				default:
					t.Errorf("ParseValue() parsed type = %T, want JSON type", parsed)
				}

				// For "Simple JSON Object" specifically:
				if tt.name == "Simple JSON Object" {
					m, ok := parsed.(map[string]any)
					if !ok || m["foo"] != "bar" {
						t.Errorf("Failed to parse JSON object correctly: %v", parsed)
					}
				}
			} else {
				// For plain text, parsed should equal raw
				if parsed != raw {
					t.Errorf("ParseValue() parsed = %v, want %v", parsed, raw)
				}
			}
		})
	}
}
