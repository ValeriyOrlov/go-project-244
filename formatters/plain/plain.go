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
		// Для составных типов (map, slice, struct) возвращаем [complex value]
		return "[complex value]"
	}
}

func Plain(data []gendiff.KeyCharacteristics, path []string) string {
	result := []string{}
	keyStatuses := map[string]string{
		"added":   " was added with value: ",
		"deleted": " was removed",
		"changed": " was updated. ",
	}
	propertyStr := "Property "

	for i, key := range data {
		var addedValue, removedValue any
		var nextKey, prevKey string
		if i > 0 {
			prevKey = data[i-1].Name
		}
		if i < len(data)-1 {
			nextKey = data[i+1].Name
		}

		switch {
		case key.Status == "changed":
			newPath := append(append([]string{}, path...), key.Name, ".")
			result = append(result, Plain(key.Value.([]gendiff.KeyCharacteristics), newPath))
		case key.Status == "deleted" && key.Name != nextKey:
			result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), key.Name), keyStatuses[key.Status], "\n")

		case key.Status == "added" && key.Name != prevKey:
			addedValue = normalizeValue(key.Value)
			result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), key.Name), keyStatuses[key.Status], addedValue.(string), "\n")

		case key.Status == "added" && key.Name == prevKey:
			removedValue = normalizeValue(data[i-1].Value)
			addedValue = normalizeValue(key.Value)
			changedValues := fmt.Sprintf("From %s to %s", removedValue, addedValue)
			result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), key.Name), keyStatuses["changed"], changedValues, "\n")
		}
	}
	return strings.Join(result, "")
}
