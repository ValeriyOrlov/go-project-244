package formatters

import (
	"code/cmd/gendiff"
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	tests := []struct {
		name     string
		data     []gendiff.KeyCharacteristics
		expected string
	}{
		{
			name:     "empty slice",
			data:     []gendiff.KeyCharacteristics{},
			expected: `[]`,
		},
		{
			name: "single added item",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Value: "value1", Status: "added"},
			},
			expected: `[{"Name":"key1","Value":"value1","Status":"added"}]`,
		},
		{
			name: "single deleted item with number",
			data: []gendiff.KeyCharacteristics{
				{Name: "key2", Value: 123, Status: "deleted"},
			},
			expected: `[{"Name":"key2","Value":123,"Status":"deleted"}]`,
		},
		{
			name: "multiple items with different types",
			data: []gendiff.KeyCharacteristics{
				{Name: "a", Value: true, Status: "added"},
				{Name: "b", Value: nil, Status: "deleted"},
				{Name: "c", Value: 3.14, Status: "equal"},
			},
			expected: `[{"Name":"a","Value":true,"Status":"added"},{"Name":"b","Value":null,"Status":"deleted"},{"Name":"c","Value":3.14,"Status":"equal"}]`,
		},
		{
			name: "nested changed item",
			data: []gendiff.KeyCharacteristics{
				{
					Name: "nested",
					Value: []gendiff.KeyCharacteristics{
						{Name: "sub", Value: "old", Status: "deleted"},
						{Name: "sub", Value: "new", Status: "added"},
					},
					Status: "changed",
				},
			},
			expected: `[{"Name":"nested","Value":[{"Name":"sub","Value":"old","Status":"deleted"},{"Name":"sub","Value":"new","Status":"added"}],"Status":"changed"}]`,
		},
		{
			name: "complex value as map",
			data: []gendiff.KeyCharacteristics{
				{Name: "mapKey", Value: map[string]interface{}{"inner": "value"}, Status: "added"},
			},
			expected: `[{"Name":"mapKey","Value":{"inner":"value"},"Status":"added"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Json(tt.data)
			if got != tt.expected {
				t.Errorf("Json() = %v, want %v", got, tt.expected)
			}
			// Дополнительная проверка: убеждаемся, что результат — валидный JSON
			var parsed interface{}
			err := json.Unmarshal([]byte(got), &parsed)
			if err != nil {
				t.Errorf("Json() produced invalid JSON: %v", err)
			}
		})
	}
}
