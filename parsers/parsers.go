package parsers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func fileReader(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
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
			return nil, fmt.Errorf("cannot read file: %w", err)
		}
		return res, nil
	}
	return nil, nil
}

func jsonParser(data []byte) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("cannot parse file: %w", err)
	}
	return result, nil
}
