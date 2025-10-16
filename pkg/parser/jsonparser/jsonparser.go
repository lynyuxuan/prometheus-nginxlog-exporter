package jsonparser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JsonParser parse a JSON string.
type JsonParser struct{}

// NewJsonParser returns a new json parser.
func NewJsonParser() *JsonParser {
	return &JsonParser{}
}

// ParseString implements the Parser interface.
// The value in the map is not necessarily a string, so it needs to be converted.
func (j *JsonParser) ParseString(line string) (map[string]string, error) {
	var parsed map[string]interface{}

	// First, try to unmarshal the whole line (fast path for pure-JSON lines)
	err := json.Unmarshal([]byte(line), &parsed)
	if err == nil {
		return toStringMap(parsed), nil
	}

	// If full-line unmarshal failed, try to extract a JSON substring.
	// This handles lines like: "<timestamp> stdout F { ... }"
	l := strings.TrimSpace(line)
	// find first '{' and last '}'
	start := strings.Index(l, "{")
	end := strings.LastIndex(l, "}")
	if start != -1 && end != -1 && end > start {
		candidate := l[start : end+1]
		// Try to unmarshal the candidate into a fresh map
		var fallback map[string]interface{}
		if uerr := json.Unmarshal([]byte(candidate), &fallback); uerr == nil {
			return toStringMap(fallback), nil
		} else {
			// combine errors for better diagnostics
			err = fmt.Errorf("%v; fallback candidate parse err: %w", err, uerr)
		}
	}

	return nil, fmt.Errorf("json log parsing err: %w", err)
}

func toStringMap(parsed map[string]interface{}) map[string]string {
	fields := make(map[string]string, len(parsed))
	for k, v := range parsed {
		if s, ok := v.(string); ok {
			fields[k] = s
		} else {
			fields[k] = fmt.Sprintf("%v", v)
		}
	}
	return fields
}
