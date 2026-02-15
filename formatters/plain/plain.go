package formatters

import (
	"code/cmd/gendiff"
	"fmt"
	"strings"
)

func normalizeValue(v any) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case bool:
		return fmt.Sprintf("%v", val)
	case string:
		return fmt.Sprintf("'%s'", val)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", val)
	default:
		return "[complex value]"
	}
}

func Plain(data []gendiff.KeyCharacteristics, path []string) string {
	var rows []string

	keyStatuses := map[string]string{
		"added":   " was added with value: ",
		"deleted": " was removed",
		"changed": " was updated. ",
	}
	propertyStr := "Property "

	for i, key := range data {
		var nextKey, prevKey string
		if i > 0 {
			prevKey = data[i-1].Name
		}
		if i < len(data)-1 {
			nextKey = data[i+1].Name
		}

		fullKey := strings.Join(path, "") + key.Name

		switch {
		case key.Status == "changed":
			newPath := append(append([]string{}, path...), key.Name, ".")
			sub := Plain(key.Value.([]gendiff.KeyCharacteristics), newPath)
			if sub != "" {
				subRows := strings.Split(sub, "\n")
				rows = append(rows, subRows...)
			}
		case key.Status == "deleted" && key.Name != nextKey:
			row := fmt.Sprintf("%s'%s'%s", propertyStr, fullKey, keyStatuses[key.Status])
			rows = append(rows, row)
		case key.Status == "added" && key.Name != prevKey:
			addedValue := normalizeValue(key.Value)
			row := fmt.Sprintf("%s'%s'%s%s", propertyStr, fullKey, keyStatuses[key.Status], addedValue)
			rows = append(rows, row)
		case key.Status == "added" && key.Name == prevKey:
			removedValue := normalizeValue(data[i-1].Value)
			addedValue := normalizeValue(key.Value)
			row := fmt.Sprintf("%s'%s'%sFrom %s to %s", propertyStr, fullKey, keyStatuses["changed"], removedValue, addedValue)
			rows = append(rows, row)
		}
	}
	return strings.Join(rows, "\n")
}
