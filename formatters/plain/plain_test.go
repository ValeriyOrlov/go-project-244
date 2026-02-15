package formatters

import (
	"code/cmd/gendiff"
	"testing"
)

func TestPlain(t *testing.T) {
	tests := []struct {
		name     string
		data     []gendiff.KeyCharacteristics
		expected string
	}{
		{
			name: "added simple value",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "added", Value: "value1"},
			},
			expected: "Property 'key1' was added with value: 'value1'",
		},
		{
			name: "deleted simple value",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "deleted", Value: "value1"},
			},
			expected: "Property 'key1' was removed",
		},
		{
			name: "updated simple value (deleted+added pair)",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "deleted", Value: "old"},
				{Name: "key1", Status: "added", Value: "new"},
			},
			expected: "Property 'key1' was updated. From 'old' to 'new'",
		},
		{
			name: "updated with boolean and null",
			data: []gendiff.KeyCharacteristics{
				{Name: "flag", Status: "deleted", Value: true},
				{Name: "flag", Status: "added", Value: nil},
			},
			expected: "Property 'flag' was updated. From true to null",
		},
		{
			name: "updated with complex value",
			data: []gendiff.KeyCharacteristics{
				{Name: "obj", Status: "deleted", Value: map[string]any{"a": 1}},
				{Name: "obj", Status: "added", Value: "simple"},
			},
			expected: "Property 'obj' was updated. From [complex value] to 'simple'",
		},
		{
			name: "nested object with added property",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "common",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "setting1", Status: "added", Value: "value1"},
					},
				},
			},
			expected: "Property 'common.setting1' was added with value: 'value1'",
		},
		{
			name: "nested object with updated property",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "common",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "setting2", Status: "deleted", Value: "old"},
						{Name: "setting2", Status: "added", Value: "new"},
					},
				},
			},
			expected: "Property 'common.setting2' was updated. From 'old' to 'new'",
		},
		{
			name: "deep nested path",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "a",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{
							Name:   "b",
							Status: "changed",
							Value: []gendiff.KeyCharacteristics{
								{Name: "c", Status: "added", Value: 42},
							},
						},
					},
				},
			},
			expected: "Property 'a.b.c' was added with value: 42",
		},
		{
			name: "multiple changes at root",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "deleted", Value: 1},
				{Name: "key1", Status: "added", Value: 2},
				{Name: "key2", Status: "added", Value: "new"},
				{Name: "key3", Status: "deleted", Value: "old"},
			},
			expected: "Property 'key1' was updated. From 1 to 2\n" +
				"Property 'key2' was added with value: 'new'\n" +
				"Property 'key3' was removed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Plain(tt.data, []string{})
			if result != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, result)
			}
		})
	}
}
