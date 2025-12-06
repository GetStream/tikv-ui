package utils

import "strings"

// SplitAndTrim splits a string by separator and trims whitespace from each part
func SplitAndTrim(s, sep string) []string {
	raw := strings.Split(s, sep)
	out := make([]string, 0, len(raw))
	for _, part := range raw {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
