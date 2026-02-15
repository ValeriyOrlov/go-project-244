package code

import "code/formatters"

// GenDiff сравнивает два файла конфигурации и возвращает разницу в заданном формате.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	return formatters.Format(filepath1, filepath2, format)
}
