package formatters

import (
	"code/cmd/gendiff"
	Json "code/formatters/json"
	Plain "code/formatters/plain"
	Stylish "code/formatters/stylish"
	"code/parsers"
	"os"
	"path/filepath"
	"testing"
)

// testFile создаёт временный файл с заданным содержимым и возвращает его путь.
func testFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestFormatSuccess(t *testing.T) {
	tmpDir := t.TempDir()

	// Содержимое первого файла (JSON)
	file1Content := `{
		"common": {
			"setting1": "Value 1",
			"setting2": 200,
			"setting3": true,
			"setting6": {
				"key": "value",
				"doge": {
					"wow": ""
				}
			}
		},
		"group1": {
			"baz": "bas",
			"foo": "bar",
			"nest": {
				"key": "value"
			}
		},
		"group2": {
			"abc": 12345,
			"deep": {
				"id": 45
			}
		}
	}`

	// Содержимое второго файла (JSON)
	file2Content := `{
		"common": {
			"follow": false,
			"setting1": "Value 1",
			"setting3": null,
			"setting4": "blah blah",
			"setting5": {
				"key5": "value5"
			},
			"setting6": {
				"key": "value",
				"ops": "vops",
				"doge": {
					"wow": "so much"
				}
			}
		},
		"group1": {
			"foo": "bar",
			"baz": "bars",
			"nest": "str"
		},
		"group3": {
			"deep": {
				"id": {
					"number": 45
				}
			},
			"fee": 100500
		}
	}`

	path1 := testFile(t, tmpDir, "file1.json", file1Content)
	path2 := testFile(t, tmpDir, "file2.json", file2Content)

	// Получаем данные и diff для сравнения с ожидаемыми результатами
	data1, err := parsers.Parser(path1)
	if err != nil {
		t.Fatal(err)
	}
	data2, err := parsers.Parser(path2)
	if err != nil {
		t.Fatal(err)
	}
	diff := gendiff.Gendiff(data1, data2)

	// Ожидаемые результаты для каждого формата
	expectedPlain := Plain.Plain(diff, nil)
	expectedJSON := Json.Json(diff)
	expectedStylish := Stylish.Stylish(diff, 0)

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{"plain", "plain", expectedPlain},
		{"json", "json", expectedJSON},
		{"stylish", "stylish", expectedStylish},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Format(path1, path2, tt.format)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Format() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	path1 := testFile(t, tmpDir, "file1.json", `{"a":1}`)
	path2 := filepath.Join(tmpDir, "nonexistent.json") // не существует

	_, err := Format(path1, path2, "plain")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestFormatInvalidFileContent(t *testing.T) {
	tmpDir := t.TempDir()
	// Создаём файл с невалидным JSON
	path1 := testFile(t, tmpDir, "bad.json", `{"a":1,}`)
	path2 := testFile(t, tmpDir, "good.json", `{"b":2}`)

	_, err := Format(path1, path2, "plain")
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestFormatUnknownFormat(t *testing.T) {
	tmpDir := t.TempDir()
	path1 := testFile(t, tmpDir, "file1.json", `{"a":1}`)
	path2 := testFile(t, tmpDir, "file2.json", `{"b":2}`)

	_, err := Format(path1, path2, "unknown")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
	if err != nil && err.Error() != "unknown format: unknown" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFormatWithYAMLFiles(t *testing.T) {
	tmpDir := t.TempDir()

	file1Content := `
common:
  setting1: Value 1
  setting2: 200
  setting3: true
  setting6:
    key: value
    doge:
      wow: ""
group1:
  baz: bas
  foo: bar
  nest:
    key: value
group2:
  abc: 12345
  deep:
    id: 45
`

	file2Content := `
common:
  follow: false
  setting1: Value 1
  setting3: null
  setting4: blah blah
  setting5:
    key5: value5
  setting6:
    key: value
    ops: vops
    doge:
      wow: so much
group1:
  foo: bar
  baz: bars
  nest: str
group3:
  deep:
    id:
      number: 45
  fee: 100500
`

	path1 := testFile(t, tmpDir, "file1.yaml", file1Content)
	path2 := testFile(t, tmpDir, "file2.yaml", file2Content)

	// Вычисляем ожидаемый результат для plain, используя данные из файлов
	data1, err := parsers.Parser(path1)
	if err != nil {
		t.Fatal(err)
	}
	data2, err := parsers.Parser(path2)
	if err != nil {
		t.Fatal(err)
	}
	diff := gendiff.Gendiff(data1, data2)
	expected := Plain.Plain(diff, nil)

	result, err := Format(path1, path2, "plain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != expected {
		t.Errorf("Format(plain) with YAML = %q, want %q", result, expected)
	}
}
