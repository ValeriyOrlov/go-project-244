package gendiff

import (
	"reflect"
	"testing"
)

func TestGendiff(t *testing.T) {
	tests := []struct {
		name     string
		data1    map[string]any
		data2    map[string]any
		expected map[string]KeyCharacteristics
	}{
		{
			name:  "одинаковые карты",
			data1: map[string]any{"key1": "value1", "key2": 2},
			data2: map[string]any{"key1": "value1", "key2": 2},
			expected: map[string]KeyCharacteristics{
				"key1": {Key: "key1", HasDiff: false, Status: "equal", Value: "value1"},
				"key2": {Key: "key2", HasDiff: false, Status: "equal", Value: 2},
			},
		},
		{
			name:  "разные значения",
			data1: map[string]any{"key1": "value1"},
			data2: map[string]any{"key1": "value2"},
			expected: map[string]KeyCharacteristics{
				"key1": {
					Key:            "key1",
					HasDiff:        false,
					IsValueChanged: true,
					Status:         "changed",
					Value: []KeyCharacteristics{
						{Key: "key1", HasDiff: false, Status: "deleted", Value: "value1"},
						{Key: "key1", HasDiff: false, Status: "added", Value: "value2"},
					},
				},
			},
		},
		{
			name:  "вложенные карты",
			data1: map[string]any{"nested": map[string]any{"inner": 1}},
			data2: map[string]any{"nested": map[string]any{"inner": 2}},
			expected: map[string]KeyCharacteristics{
				"nested": {
					Key:            "nested",
					HasDiff:        true,
					IsValueChanged: false,
					Status:         "diff",
					Value: map[string]KeyCharacteristics{
						"inner": {
							Key:            "inner",
							HasDiff:        false,
							IsValueChanged: true,
							Status:         "changed",
							Value: []KeyCharacteristics{
								{Key: "inner", HasDiff: false, Status: "deleted", Value: 1},
								{Key: "inner", HasDiff: false, Status: "added", Value: 2},
							},
						},
					},
				},
			},
		},
		{
			name:  "ключ только в data1",
			data1: map[string]any{"onlyInFirst": 123},
			data2: map[string]any{},
			expected: map[string]KeyCharacteristics{
				"onlyInFirst": {Key: "onlyInFirst", HasDiff: false, Status: "deleted", Value: 123},
			},
		},
		{
			name:  "ключ только в data2",
			data1: map[string]any{},
			data2: map[string]any{"onlyInSecond": true},
			expected: map[string]KeyCharacteristics{
				"onlyInSecond": {Key: "onlyInSecond", HasDiff: false, Status: "added", Value: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Gendiff(tt.data1, tt.data2)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unexpected result for %s. Got\n %+v, want\n %+v", tt.name, result, tt.expected)
			}
		})
	}
}
