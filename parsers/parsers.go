package parsers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var errReadFile = "cannot read file: %w"
var errParseFile = "cannot parse file: %w"

func fileReader(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(errReadFile, err)
	}
	return data, nil
}

func Parser(fp string) (map[string]any, error) {
	f, err := fileReader(fp)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasSuffix(fp, "json") {
		res, err := jsonParser(f)
		if err != nil {
			return nil, fmt.Errorf(errReadFile, err)
		}
		return res, nil
	}
	if strings.HasSuffix(fp, "yml") || strings.HasSuffix(fp, "yaml") {
		res, err := ymlParser(f)
		if err != nil {
			return nil, fmt.Errorf(errReadFile, err)
		}
		return res, nil
	}
	return nil, nil
}

func jsonParser(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf(errParseFile, err)
	}
	return result, nil
}

func ymlParser(data []byte) (map[string]any, error) {
	var result map[string]any
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf(errParseFile, err)
	}
	return result, nil
}
