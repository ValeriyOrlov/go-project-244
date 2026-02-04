package formatters

import (
	"code/cmd/gendiff"
	"code/parsers"
	"fmt"
	"log"
	"maps"
	"slices"
	"strings"
)

// функция для отрисовки отступов
func marginsCreator(nestingCounter int) string {
	margin := make([]string, 0, nestingCounter)
	for i := 0; i < nestingCounter; i++ {
		margin = append(margin, "..")
	}
	return strings.Join(margin, "")
}

// функции для отрисовки плоских строк и строк с вложенными ключами
func plainRowCreator(margin, status, key string, value any) string {
	return fmt.Sprintf("%s%s %s: %v\n", margin, status, key, value)
}

func nestedRowCreator(margin, status, key string, value any) string {
	return fmt.Sprintf("%s%s %s: {\n%v", margin, status, key, value)
}

// функция для отрисовки ключей с вложенными ключами
func mapPrint(row map[string]any, nestingCounter int) string {
	result := []string{}
	margin := marginsCreator(nestingCounter)
	for key, value := range row {
		if gendiff.IsMap(value) {
			nestingCounter += 1
			newRow := nestedRowCreator(margin, "", key, mapPrint(value.(map[string]any), nestingCounter))
			result = append(result, newRow)
			closeBracket := fmt.Sprintf("%s}\n", margin)
			result = append(result, closeBracket)
			nestingCounter -= 1
		} else {
			newRow := plainRowCreator(margin, "", key, value)
			result = append(result, newRow)
		}
	}
	return strings.Join(result, "")
}

func Stylish(data map[string]gendiff.KeyCharacteristics, nestingCounter int) string {
	result := make([]string, len(data))
	if nestingCounter == 0 {
		result = append(result, "{\n")
		nestingCounter += 1
	}
	//достаем и сортируем ключи
	keys := slices.Sorted(maps.Keys(data))
	//вводим карту со статусами для строк, чтобы при необходимости поменять знак характеризующий изменения можно было в одном месте
	keyStatuses := map[string]string{"added": "+", "deleted": "-", "equal": " ", "diff": " "}
	margin := marginsCreator(nestingCounter)
	closeBracket := fmt.Sprintf("%s}\n", margin)
	//проходим по срезу с ключами
	for _, key := range keys {
		//если ключ содержит дифф, строим строку с ключем и для отрисовки вложенных ключей запускаем форматтер рекурсивно
		switch {
		case data[key].HasDiff:
			nestingCounter += 1
			row := fmt.Sprintf("%s%s %v: {\n", margin, keyStatuses[data[key].Status], data[key].Key)
			result = append(result, row)
			result = append(result, Stylish(data[key].Value.(map[string]gendiff.KeyCharacteristics), nestingCounter))
			result = append(result, closeBracket)
			nestingCounter -= 1
		//здесь рекурсивно отрисовывем карты
		case gendiff.IsMap(data[key].Value):
			nestingCounter += 1
			row := nestedRowCreator(margin, keyStatuses[data[key].Status], data[key].Key, mapPrint(data[key].Value.(map[string]any), nestingCounter))
			result = append(result, row)
			result = append(result, closeBracket)
			nestingCounter -= 1
		//здесь отрисовываются случаи, когда в обоих файлах найдены одинаковые ключи с разными значениями. это единственный случай, где значение ключа является срезом, поэтому для корректной работы проверяем,является ли ключ срезом
		case data[key].IsValueChanged && gendiff.IsSlice(data[key].Value):
			changedValues := data[key].Value.([]gendiff.KeyCharacteristics)
			for _, val := range changedValues {
				if gendiff.IsMap(val.Value) {
					nestingCounter += 1
					row := nestedRowCreator(margin, keyStatuses[val.Status], val.Key, mapPrint((val.Value.(map[string]any)), nestingCounter))
					result = append(result, row)
					result = append(result, closeBracket)
					nestingCounter -= 1
				} else {
					row := plainRowCreator(margin, keyStatuses[val.Status], val.Key, val.Value)
					result = append(result, row)
				}
			}
		//здесь отрисовываются ключи со значениями без вложенных объектов
		case !gendiff.IsMap(data[key].Value):
			row := plainRowCreator(margin, keyStatuses[data[key].Status], data[key].Key, data[key].Value)
			result = append(result, row)
		}
	}
	if nestingCounter == 1 {
		result = append(result, "}")
	}
	return fmt.Sprint(strings.Join(result, ""))
}

func Formatters(path1 string, path2 string, stylish bool) string {
	data1, err1 := parsers.Parser(path1)
	if err1 != nil {
		log.Fatal(err1)
	}
	data2, err2 := parsers.Parser(path2)
	if err2 != nil {
		log.Fatal(err2)
	}
	diff := gendiff.Gendiff(data1, data2)
	if stylish {
		return Stylish(diff, 0)
	}
	return ""
}
