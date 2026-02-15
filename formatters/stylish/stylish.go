package formatters

import (
	"code/cmd/gendiff"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// indent возвращает отступ из level * 4 пробелов
func indent(level int) string {
	return strings.Repeat("    ", level)
}

// stringify преобразует значение в строку, заменяя nil на "null", а строки оборачивает в кавычки
func stringify(v any) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// renderObject форматирует map как многострочный объект с отступами
func renderObject(obj map[string]any, level int) string {
	if len(obj) == 0 {
		return "{}"
	}
	var builder strings.Builder
	builder.WriteString("{\n")
	keys := slices.Sorted(maps.Keys(obj))
	for _, k := range keys {
		v := obj[k]
		builder.WriteString(indent(level + 1))
		if gendiff.IsMap(v) {
			builder.WriteString(fmt.Sprintf("%s: %s", k, renderObject(v.(map[string]any), level+1)))
		} else {
			builder.WriteString(fmt.Sprintf("%s: %s", k, stringify(v)))
		}
		builder.WriteString("\n")
	}
	builder.WriteString(indent(level) + "}")
	return builder.String()
}

// Stylish возвращает отформатированный diff в стиле stylish
func Stylish(data []gendiff.KeyCharacteristics, level int) string {
	statusSign := map[string]string{
		"added":   "+",
		"deleted": "-",
		"equal":   " ",
		"changed": " ",
	}

	var body strings.Builder

	for _, item := range data {
		sign := statusSign[item.Status]
		body.WriteString(indent(level + 1))

		switch item.Status {
		case "changed":
			// Вложенные изменения
			body.WriteString(fmt.Sprintf("  %s: ", item.Name))
			// Рекурсивно получаем форматирование вложенного diff
			nested := Stylish(item.Value.([]gendiff.KeyCharacteristics), level+1)
			body.WriteString(nested)
		default:
			if gendiff.IsMap(item.Value) {
				body.WriteString(fmt.Sprintf("%s %s: %s", sign, item.Name, renderObject(item.Value.(map[string]any), level+1)))
			} else {
				body.WriteString(fmt.Sprintf("%s %s: %s", sign, item.Name, stringify(item.Value)))
			}
			body.WriteString("\n")
		}
	}

	// Оборачиваем тело в фигурные скобки с правильными отступами
	result := strings.Builder{}
	result.WriteString("{\n")
	result.WriteString(body.String())
	result.WriteString(indent(level) + "}\n")
	return result.String()
}
