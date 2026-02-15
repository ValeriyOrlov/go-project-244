package formatters

import (
	"code/cmd/gendiff"
	Json "code/formatters/json"
	Plain "code/formatters/plain"
	Stylish "code/formatters/stylish"
	"code/parsers"
	"fmt"
)

func Format(path1, path2, format string) (string, error) {
	data1, err := parsers.Parser(path1)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", path1, err)
	}
	data2, err := parsers.Parser(path2)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", path2, err)
	}
	diff := gendiff.Gendiff(data1, data2)
	switch format {
	case "plain":
		return Plain.Plain(diff, nil), nil
	case "json":
		return Json.Json(diff), nil
	case "stylish":
		return Stylish.Stylish(diff, 0), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
