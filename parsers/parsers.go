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

func Parser(fp string) {
	f, err := fileReader(fp)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasSuffix(fp, "json") {
		jsonParser(f)
	}
}

func jsonParser(data []byte) error {
	var result any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("cannot parse file: %w", err)
	}
	fmt.Printf("%+v\n", result)
	return nil
}
