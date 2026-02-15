package formatters

import (
	"code/cmd/gendiff"
	"testing"
)

func TestStylish(t *testing.T) {
	tests := []struct {
		name     string
		data     []gendiff.KeyCharacteristics
		expected string
	}{
		{
			name:     "empty diff",
			data:     []gendiff.KeyCharacteristics{},
			expected: "{\n}\n",
		},
		{
			name: "single added key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Value: "value1", Status: "added"},
			},
			expected: "{\n    + key1: 'value1'\n}\n",
		},
		{
			name: "single deleted key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Value: 42, Status: "deleted"},
			},
			expected: "{\n    - key1: 42\n}\n",
		},
		{
			name: "single unchanged key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "key1", Value: true, Status: "equal"},
			},
			expected: "{\n      key1: true\n}\n",
		},
		{
			name: "multiple keys with all statuses – sorted",
			data: []gendiff.KeyCharacteristics{
				{Name: "a", Value: 1, Status: "added"},
				{Name: "b", Value: 2, Status: "deleted"},
				{Name: "c", Value: 3, Status: "equal"},
			},
			expected: "{\n    + a: 1\n    - b: 2\n      c: 3\n}\n",
		},
		{
			name: "changed – nested equal diff, level 2 inside",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "nested",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "inner", Value: "value", Status: "equal"},
					},
				},
			},
			expected: "{\n      nested: {\n          inner: 'value'\n    }\n}\n",
		},
		{
			name: "changed – nested diff with changes, level 2",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "nested",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "inner", Value: "old", Status: "deleted"},
						{Name: "inner", Value: "new", Status: "added"},
					},
				},
			},
			expected: "{\n      nested: {\n        - inner: 'old'\n        + inner: 'new'\n    }\n}\n",
		},
		{
			name: "map value added as whole object – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "obj",
					Value:  map[string]any{"b": 2, "a": 1},
					Status: "added",
				},
			},
			expected: "{\n    + obj: {\n        a: 1\n        b: 2\n    }\n}\n",
		},
		{
			name: "map value deleted as whole object – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "obj",
					Value:  map[string]any{"foo": "bar", "num": 42},
					Status: "deleted",
				},
			},
			expected: "{\n    - obj: {\n        foo: 'bar'\n        num: 42\n    }\n}\n",
		},
		{
			name: "unchanged map value (equal) – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "obj",
					Value:  map[string]any{"x": 1, "y": 2},
					Status: "equal",
				},
			},
			expected: "{\n      obj: {\n        x: 1\n        y: 2\n    }\n}\n",
		},
		{
			name: "null values",
			data: []gendiff.KeyCharacteristics{
				{Name: "nullKey", Value: nil, Status: "added"},
				{Name: "nullKey", Value: nil, Status: "deleted"},
			},
			expected: "{\n    + nullKey: null\n    - nullKey: null\n}\n",
		},
		{
			name: "deep nested mix – all rules together",
			data: []gendiff.KeyCharacteristics{
				{Name: "added", Value: "new", Status: "added"},
				{Name: "deleted", Value: 123, Status: "deleted"},
				{Name: "equal", Value: "same", Status: "equal"},
				{
					Name:   "changed",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "innerAdded", Value: "inner", Status: "added"},
						{Name: "innerDeleted", Value: "old", Status: "deleted"},
						{Name: "innerEqual", Value: "same", Status: "equal"},
						{
							Name:   "innerChanged",
							Status: "changed",
							Value: []gendiff.KeyCharacteristics{
								{Name: "deep", Value: "deepValue", Status: "added"},
							},
						},
					},
				},
				{
					Name:   "objAdded",
					Value:  map[string]any{"a": 1, "b": 2},
					Status: "added",
				},
				{
					Name:   "objDeleted",
					Value:  map[string]any{"c": 3, "d": 4},
					Status: "deleted",
				},
				{
					Name:   "objEqual",
					Value:  map[string]any{"e": 5, "f": 6},
					Status: "equal",
				},
			},
			expected: `{
    + added: 'new'
    - deleted: 123
      equal: 'same'
      changed: {
        + innerAdded: 'inner'
        - innerDeleted: 'old'
          innerEqual: 'same'
          innerChanged: {
            + deep: 'deepValue'
        }
    }
    + objAdded: {
        a: 1
        b: 2
    }
    - objDeleted: {
        c: 3
        d: 4
    }
      objEqual: {
        e: 5
        f: 6
    }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Stylish(tt.data, 0)
			if result != tt.expected {
				t.Errorf("Stylish() = %q, want %q", result, tt.expected)
			}
		})
	}
}
