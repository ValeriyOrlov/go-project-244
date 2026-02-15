package parsers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parser(fp string) (map[string]any, error) {
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %w", fp, err)
	}

	switch {
	case strings.HasSuffix(fp, ".json"):
		return parseJSON(data)
	case strings.HasSuffix(fp, ".yml"), strings.HasSuffix(fp, ".yaml"):
		return parseYAML(data)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", fp)
	}
}

func parseJSON(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("cannot parse JSON: %w", err)
	}
	return result, nil
}

func parseYAML(data []byte) (map[string]any, error) {
	var result map[string]any
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("cannot parse YAML: %w", err)
	}
	return result, nil
}
