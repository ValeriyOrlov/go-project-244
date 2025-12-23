package formatters

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func Stylish(data map[string]string) string {
	sortedDiffKeys := slices.Sorted(maps.Keys(data))
	result := make([]string, len(data)+2)
	result = append(result, "{\n")
	for _, key := range sortedDiffKeys {
		row := fmt.Sprintln(" ", data[key], key)
		result = append(result, row)
	}
	result = append(result, "}\n")
	return strings.Join(result, "")
}
