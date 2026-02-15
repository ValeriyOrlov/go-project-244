package formatters

import (
	"code/cmd/gendiff"
	"testing"
)

func TestStylish(t *testing.T) {
	tests := []struct {
		name     string
		data     []gendiff.KeyCharacteristics
		level    int
		expected string
	}{
		{
			name:     "empty diff",
			data:     []gendiff.KeyCharacteristics{},
			level:    0,
			expected: "{\n}",
		},
		{
			name: "single equal key",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "equal", Value: "value1"},
			},
			level:    0,
			expected: "{\n    key1: value1\n}",
		},
		{
			name: "added and deleted keys",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Status: "added", Value: "new value"},
				{Name: "key2", Status: "deleted", Value: "old value"},
			},
			level:    0,
			expected: "{\n  + key1: new value\n  - key2: old value\n}",
		},
		{
			name: "nested map (equal)",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "nested",
					Status: "equal",
					Value: map[string]any{
						"a": 1,
						"b": "two",
					},
				},
			},
			level:    0,
			expected: "{\n    nested: {\n        a: 1\n        b: two\n    }\n}",
		},
		{
			name: "changed status with nested diff",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "changedKey",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "sub1", Status: "added", Value: "subValue1"},
						{Name: "sub2", Status: "deleted", Value: "subValue2"},
					},
				},
			},
			level:    0,
			expected: "{\n    changedKey: {\n      + sub1: subValue1\n      - sub2: subValue2\n    }\n}",
		},
		{
			name: "mixed keys with different statuses",
			data: []gendiff.KeyCharacteristics{
				{Name: "common", Status: "equal", Value: "same"},
				{Name: "toRemove", Status: "deleted", Value: 42},
				{Name: "toAdd", Status: "added", Value: true},
				{
					Name:   "nestedChanged",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "inside", Status: "added", Value: "inner"},
					},
				},
				{
					Name:   "object",
					Status: "equal",
					Value: map[string]any{
						"x": nil,
						"y": "str",
					},
				},
			},
			level: 0,
			expected: "{\n" +
				"    common: same\n" +
				"  - toRemove: 42\n" +
				"  + toAdd: true\n" +
				"    nestedChanged: {\n      + inside: inner\n    }\n" +
				"    object: {\n        x: null\n        y: str\n    }\n" +
				"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Stylish(tt.data, tt.level)
			if got != tt.expected {
				t.Errorf("\ngot:\n %q\n want:\n %q\n", got, tt.expected)
			}
		})
	}
}
