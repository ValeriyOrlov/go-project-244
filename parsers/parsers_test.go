package parsers

import (
	"maps"
	"slices"
	"testing"
)

func TestFileReaderSuccess(t *testing.T) {
	data, err := fileReader("../testdata/fixtures/file1.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected data to be non-empty")
	}
}

func TestFileReaderError(t *testing.T) {
	_, err := fileReader("../testdata/fixtures/nonexistentfile.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestJsonParserValid(t *testing.T) {
	jsonData := []byte(`{"key": "value"}`)
	result, err := jsonParser(jsonData)
	if err != nil {
		t.Fatal(err)
	}
	if result["key"] != "value" {
		t.Fatal("parsed result does not contain expected key-value")
	}
}

func TestJsonParserInvalid(t *testing.T) {
	jsonData := []byte(`{key: value}`)
	_, err := jsonParser(jsonData)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParserJSONFile(t *testing.T) {
	result, err := Parser("../testdata/fixtures/file1.json")
	if err != nil {
		t.Fatal(err)
	}
	gotKeys := slices.Sorted(maps.Keys(result))
	wantKeys := []string{"follow", "host", "proxy", "timeout"}
	if !slices.Equal(gotKeys, wantKeys) {
		t.Fatal(err)
	}
}
