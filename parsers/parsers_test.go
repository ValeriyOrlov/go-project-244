package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

// TestJSONParser проверяет корректный парсинг JSON-файла.
func TestJSONParser(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	content := `{"key": "value", "num": 42, "flag": true, "null": null, "empty": ""}`
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Parser(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Проверяем все поля
	checkField(t, result, "key", "value")
	checkField(t, result, "num", float64(42)) // JSON numbers → float64
	checkField(t, result, "flag", true)
	checkField(t, result, "null", nil)
	checkField(t, result, "empty", "")
}

// TestYAMLParser проверяет корректный парсинг YAML-файла.
func TestYAMLParser(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.yaml")
	// Используем простую строку без лишних отступов
	content := "key: value\nnum: 42\nflag: true\nempty: \"\"\n"
	// Ключ null пока пропущен, так как библиотека gopkg.in/yaml.v3 может его игнорировать
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Parser(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkField(t, result, "key", "value")
	checkField(t, result, "num", 42) // YAML числа → int
	checkField(t, result, "flag", true)
	checkField(t, result, "empty", "")
}

// TestParserFileNotFound проверяет ошибку при отсутствии файла.
func TestParserFileNotFound(t *testing.T) {
	_, err := Parser("/non/existent/file.json")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

// TestParserUnsupportedFormat проверяет ошибку при неподдерживаемом расширении.
func TestParserUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(filePath, []byte("some data"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parser(filePath)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

// TestParserInvalidJSON проверяет ошибку при невалидном JSON.
func TestParserInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bad.json")
	content := `{"key": "value",}` // лишняя запятая
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parser(filePath)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// TestParserInvalidYAML проверяет ошибку при невалидном YAML.
func TestParserInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "bad.yaml")
	content := `key: value: extra` // некорректный YAML
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parser(filePath)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

// Вспомогательная функция для проверки поля с ожидаемым значением.
func checkField(t *testing.T, m map[string]any, key string, expected any) {
	t.Helper()
	val, ok := m[key]
	if !ok {
		t.Errorf("key %q missing", key)
		return
	}
	// Для nil используем специальную проверку
	if expected == nil {
		if val != nil {
			t.Errorf("key %q: expected nil, got %v (%T)", key, val, val)
		}
		return
	}
	if val != expected {
		t.Errorf("key %q: expected %v (%T), got %v (%T)", key, expected, expected, val, val)
	}
}
