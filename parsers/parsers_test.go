package parsers

import (
	"reflect"
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

func TestJsonParser(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]any
		wantErr bool
	}{
		{
			name:  "валидный вложенный JSON",
			input: []byte(`{"common": {"setting1": "Value 1"}}`),
			want:  map[string]any{"common": map[string]any{"setting1": "Value 1"}},
		},
		{
			name:    "некорректный JSON",
			input:   []byte(`{"key": "value"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonParser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("jsonParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jsonParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYmlParser(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]any
		wantErr bool
	}{

		{
			name: "валидный вложенный YAML",
			input: []byte(`
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
`),
			want: map[string]any{
				"common": map[string]any{
					"follow":   false,
					"setting1": "Value 1",
					"setting3": nil,
					"setting4": "blah blah",
					"setting5": map[string]any{"key5": "value5"},
					"setting6": map[string]any{"key": "value", "ops": "vops", "doge": map[string]any{"wow": "so much"}},
				},
				"group1": map[string]any{
					"foo":  "bar",
					"baz":  "bars",
					"nest": "str",
				},
				"group3": map[string]any{
					"deep": map[string]any{
						"id": map[string]any{
							"number": 45,
						},
					},
					"fee": 100500,
				},
			},
		},
		{
			name:    "некорректный YAML",
			input:   []byte(`: invalid yaml`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ymlParser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ymlParser() error =\n %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ymlParser() =\n %v, \nwant %v", got, tt.want)
			}
		})
	}
}
