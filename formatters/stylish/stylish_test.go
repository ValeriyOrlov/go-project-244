package formatters_test

import (
	"code/cmd/gendiff"
	formatters "code/formatters/stylish"
	"testing"
)

// Ожидаемый формат stylish в соответствии с условиями:
// - Каждый уровень вложенности добавляет 4 символа отступа.
// - Строки со статусом added/deleted: (уровень*4-2) пробела, '+'/'-', пробел, ключ: значение
// - Строки со статусом equal/changed: уровень*4 пробелов, ключ: { (для changed) или ключ: значение
// - Закрывающие скобки: уровень*4 пробелов, затем '}'
// - Внутри карт (без статуса): (уровень+1)*4 пробелов, ключ: значение (или ключ: { для вложенной карты)
// - Ключи в картах выводятся в алфавитном порядке.

func TestStylish(t *testing.T) {
	tests := []struct {
		name           string
		data           []gendiff.KeyCharacteristics
		nestingCounter int
		expected       string
	}{
		{
			name:           "empty diff",
			data:           []gendiff.KeyCharacteristics{},
			nestingCounter: 0,
			expected: `{
}
`,
		},
		{
			name: "single added key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "host", Value: "hexlet.io", Status: "added"},
			},
			nestingCounter: 0,
			expected: `{
  + host: hexlet.io
}
`,
		},
		{
			name: "single deleted key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "timeout", Value: 50, Status: "deleted"},
			},
			nestingCounter: 0,
			expected: `{
  - timeout: 50
}
`,
		},
		{
			name: "single unchanged key, level 1",
			data: []gendiff.KeyCharacteristics{
				{Name: "verbose", Value: true, Status: "equal"},
			},
			nestingCounter: 0,
			expected: `{
    verbose: true
}
`,
		},
		{
			name: "multiple keys with all statuses – sorted",
			data: []gendiff.KeyCharacteristics{
				{Name: "a", Value: 1, Status: "deleted"},
				{Name: "a", Value: 2, Status: "added"},
				{Name: "b", Value: "str", Status: "equal"},
				{Name: "c", Value: nil, Status: "added"},
			},
			nestingCounter: 0,
			expected: `{
  - a: 1
  + a: 2
    b: str
  + c: <nil>
}
`,
		},
		{
			name: "changed – nested equal diff, level 2 inside",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "group",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "x", Value: 10, Status: "equal"},
						{Name: "y", Value: 20, Status: "equal"},
					},
				},
			},
			nestingCounter: 0,
			expected: `{
    group: {
        x: 10
        y: 20
    }
}
`,
		},
		{
			name: "changed – nested diff with changes, level 2",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "group",
					Status: "changed",
					Value: []gendiff.KeyCharacteristics{
						{Name: "x", Value: 10, Status: "deleted"},
						{Name: "x", Value: 99, Status: "added"},
						{Name: "y", Value: 20, Status: "equal"},
						{Name: "z", Value: 30, Status: "added"},
					},
				},
			},
			nestingCounter: 0,
			expected: `{
    group: {
      - x: 10
      + x: 99
        y: 20
      + z: 30
    }
}
`,
		},
		{
			name: "map value added as whole object – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "proxy",
					Status: "added",
					Value: map[string]interface{}{
						"host": "localhost",
						"port": 8080,
					},
				},
			},
			nestingCounter: 0,
			expected: `{
  + proxy: {
        host: localhost
        port: 8080
    }
}
`,
		},
		{
			name: "map value deleted as whole object – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "oldConfig",
					Status: "deleted",
					Value: map[string]interface{}{
						"version": "1.0",
						"enabled": false,
					},
				},
			},
			nestingCounter: 0,
			expected: `{
  - oldConfig: {
        enabled: false
        version: 1.0
    }
}
`,
		},
		{
			name: "unchanged map value (equal) – sorted keys",
			data: []gendiff.KeyCharacteristics{
				{
					Name:   "constants",
					Status: "equal",
					Value: map[string]interface{}{
						"pi": 3.14,
						"e":  2.71,
					},
				},
			},
			nestingCounter: 0,
			expected: `{
    constants: {
        e: 2.71
        pi: 3.14
    }
}
`,
		},
		{
			name: "deep nested mix – all rules together",
			data: []gendiff.KeyCharacteristics{
				{Name: "common", Value: map[string]interface{}{"setting": "value"}, Status: "deleted"},
				{Name: "common", Value: "new value", Status: "added"},
				{Name: "group1", Status: "changed", Value: []gendiff.KeyCharacteristics{
					{Name: "a", Value: 1, Status: "equal"},
					{Name: "b", Value: 2, Status: "deleted"},
					{Name: "b", Value: 3, Status: "added"},
					{Name: "c", Status: "changed", Value: []gendiff.KeyCharacteristics{
						{Name: "inner", Value: "old", Status: "deleted"},
						{Name: "inner", Value: "new", Status: "added"},
						{Name: "extra", Value: map[string]interface{}{"flag": true}, Status: "added"},
					}},
				}},
				{Name: "group2", Value: map[string]interface{}{"x": 10, "y": 20}, Status: "added"},
				{Name: "group3", Value: "untouched", Status: "equal"},
			},
			nestingCounter: 0,
			expected: `{
  - common: {
        setting: value
    }
  + common: new value
    group1: {
        a: 1
      - b: 2
      + b: 3
        c: {
          - inner: old
          + inner: new
          + extra: {
                flag: true
            }
        }
    }
  + group2: {
        x: 10
        y: 20
    }
    group3: untouched
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatters.Stylish(tt.data, tt.nestingCounter)
			if got != tt.expected {
				t.Errorf("Stylish() output mismatch\n--- got:\n%s--- expected:\n%s", got, tt.expected)
			}
		})
	}
}
