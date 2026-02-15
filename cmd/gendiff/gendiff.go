package gendiff

import (
	"reflect"
	"sort"
)

type KeyCharacteristics struct {
	Name   string
	Value  any
	Status string
}

func IsMap(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func Gendiff(data1, data2 map[string]any) []KeyCharacteristics {
	const (
		equal   = "equal"
		deleted = "deleted"
		added   = "added"
		changed = "changed"
	)
	resultDiff := []KeyCharacteristics{}

	uniqueKeysMap := make(map[string]struct{}, len(data1)+len(data2))
	for k := range data1 {
		uniqueKeysMap[k] = struct{}{}
	}
	for k := range data2 {
		uniqueKeysMap[k] = struct{}{}
	}

	// Извлекаем ключи в срез
	uniqueKeys := make([]string, 0, len(uniqueKeysMap))
	for k := range uniqueKeysMap {
		uniqueKeys = append(uniqueKeys, k)
	}

	// Сортируем финальный срез
	sort.Strings(uniqueKeys)

	for _, key := range uniqueKeys {
		val1, ok1 := data1[key]
		val2, ok2 := data2[key]
		newKeyChar := KeyCharacteristics{}
		switch {
		//если оба значения ключа - карты, то создаем дифф (идем в рекурсию)
		case ok1 && ok2 && IsMap(val1) && IsMap(val2):
			newKeyChar.Name = key
			newKeyChar.Status = changed
			newKeyChar.Value = Gendiff(val1.(map[string]any), val2.(map[string]any))
			resultDiff = append(resultDiff, newKeyChar)
		//значения ключа не равны - добавляем в результирующий слайс ключ с измененными значениями (удаленным и добавленным)
		case ok1 && ok2 && val1 != val2:
			newKeyChar.Name = key
			newKeyChar.Status = deleted
			newKeyChar.Value = val1
			resultDiff = append(resultDiff, newKeyChar)
			newKeyChar.Status = added
			newKeyChar.Value = val2
			resultDiff = append(resultDiff, newKeyChar)
		//одинаковая пара ключ/значение
		case ok1 && ok2 && val1 == val2:
			newKeyChar.Name = key
			newKeyChar.Status = equal
			newKeyChar.Value = val1
			resultDiff = append(resultDiff, newKeyChar)
		case ok1 && !ok2:
			// ключ есть только в data1
			newKeyChar.Name = key
			newKeyChar.Status = deleted
			newKeyChar.Value = val1
			resultDiff = append(resultDiff, newKeyChar)
		case !ok1 && ok2:
			// ключ есть только в data2
			newKeyChar.Name = key
			newKeyChar.Status = added
			newKeyChar.Value = val2
			resultDiff = append(resultDiff, newKeyChar)
		}
	}
	return resultDiff
}
