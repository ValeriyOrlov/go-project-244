package formatters

import (
	"code/cmd/gendiff"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// функция для отрисовки отступов
func marginsCreator(nestingCounter int) string {
	margin := []string{}
	for i := 0; i < nestingCounter; i++ {
		margin = append(margin, "  ")
	}
	return strings.Join(margin, "")
}

// функции для отрисовки плоских строк и строк с вложенными значениями
func plainRowCreator(margin, status, key string, value any) string {
	return fmt.Sprintf("%s%s %s: %v\n", margin, status, key, value)
}

func nestedRowCreator(margin, status, key string, value any) string {
	return fmt.Sprintf("%s%s %s: %v", margin, status, key, value)
}

// функция для отрисовки карт
func mapPrint(row map[string]any, nestingCounter int) string {
	var result strings.Builder
	rowKeys := slices.Sorted(maps.Keys(row))
	nestingCounter += 1
	result.WriteString("{\n")

	for _, key := range rowKeys {
		if gendiff.IsMap(row[key]) {
			newRow := nestedRowCreator(marginsCreator(nestingCounter+1), " ", key, mapPrint(row[key].(map[string]any), nestingCounter))
			result.WriteString(newRow)
		} else {
			newRow := plainRowCreator(marginsCreator(nestingCounter+1), " ", key, row[key])
			result.WriteString(newRow)
		}
	}
	result.WriteString(marginsCreator(nestingCounter) + "}\n")

	return result.String()
}

func Stylish(data []gendiff.KeyCharacteristics, nestingCounter int) string {
	var result strings.Builder

	nestingCounter += 1
	result.WriteString("{\n")
	//карта с символами статусов
	keyStatuses := map[string]string{"added": "+", "deleted": "-", "equal": " ", "changed": " "}
	for _, key := range data {
		if key.Status == "changed" {
			row := nestedRowCreator(marginsCreator(nestingCounter), keyStatuses[key.Status], key.Name, "")
			result.WriteString(row)
			result.WriteString(Stylish(key.Value.([]gendiff.KeyCharacteristics), nestingCounter+1))
		} else if gendiff.IsMap(key.Value) {
			row := nestedRowCreator(marginsCreator(nestingCounter), keyStatuses[key.Status], key.Name, mapPrint(key.Value.(map[string]any), nestingCounter))
			result.WriteString(row)
		} else {
			row := plainRowCreator(marginsCreator(nestingCounter), keyStatuses[key.Status], key.Name, key.Value)
			result.WriteString(row)
		}
	}
	result.WriteString(marginsCreator(nestingCounter-1) + "}\n")

	return result.String()
}
