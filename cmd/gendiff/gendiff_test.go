package gendiff

import (
	"testing"
)

func TestPlainGendiff(t *testing.T) {
	tests := []struct {
		name     string
		data1    map[string]any
		data2    map[string]any
		expected map[string]string
	}{
		{
			"name: equal keys and values",
			map[string]any{"a": 1, "b": 2},
			map[string]any{"a": 1, "b": 2},
			map[string]string{
				"a: 1": " ",
				"b: 2": " ",
			},
		},
		{
			name:  "different values",
			data1: map[string]any{"a": 1},
			data2: map[string]any{"a": 2},
			expected: map[string]string{
				"a: 1": "-",
				"a: 2": "+",
			},
		},
		{
			name:  "key missing in data2",
			data1: map[string]any{"a": 1, "b": 2},
			data2: map[string]any{"a": 1},
			expected: map[string]string{
				"a: 1": " ",
				"b: 2": "-",
			},
		},
		{
			name:  "key missing in data1",
			data1: map[string]any{"a": 1},
			data2: map[string]any{"a": 1, "b": 2},
			expected: map[string]string{
				"a: 1": " ",
				"b: 2": "+",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Gendiff(tt.data1, tt.data2)

			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d entries, got %d", len(tt.expected), len(result))
			}

			for key, val := range tt.expected {
				resultVal, ok := result[key]
				if !ok {
					t.Errorf("missing key in result: %s", key)
				} else if resultVal != val {
					t.Errorf("for key %s, expected %s, got %s", key, val, resultVal)
				}
			}
		})
	}
}
