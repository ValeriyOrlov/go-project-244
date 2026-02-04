package gendiff

import (
	"maps"
	"reflect"
	"slices"
	"sort"
)

type KeyCharacteristics struct {
	Key            string
	HasDiff        bool
	IsValueChanged bool
	Value          any
	Status         string
}

func IsMap(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsSlice(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func Gendiff(data1, data2 map[string]any) map[string]KeyCharacteristics {
	const (
		equal   = "equal"
		deleted = "deleted"
		added   = "added"
		diff    = "diff"
	)
	resultDiff := make(map[string]KeyCharacteristics)
	data1Keys := slices.Sorted(maps.Keys(data1))
	data2Keys := slices.Sorted(maps.Keys(data2))
	allKeys := append(data1Keys, data2Keys...)
	// Создаем карту для хранения уникальных ключей
	uniqueMap := make(map[string]struct{})

	// Проходим по всему срезу и добавляем в карту
	for _, key := range allKeys {
		uniqueMap[key] = struct{}{}
	}
	// Формируем новый срез из ключей карты
	uniqueKeys := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		uniqueKeys = append(uniqueKeys, key)
	}

	// Опционально — сортируем для удобства
	sort.Strings(uniqueKeys)

	for _, key := range uniqueKeys {
		val1, ok1 := data1[key]
		val2, ok2 := data2[key]
		newKeyChar := KeyCharacteristics{}

		switch {
		// Оба значения - карты
		case ok1 && ok2 && IsMap(val1) && IsMap(val2):
			newKeyChar.Key = key
			newKeyChar.HasDiff = true
			newKeyChar.IsValueChanged = false
			newKeyChar.Status = diff
			newKeyChar.Value = Gendiff(val1.(map[string]any), val2.(map[string]any))
			resultDiff[key] = newKeyChar
		case ok1 && ok2 && val1 == val2:
			// Значения не карты и равны
			newKeyChar.Key = key
			newKeyChar.HasDiff = false
			newKeyChar.Status = equal
			newKeyChar.Value = val1
			resultDiff[key] = newKeyChar
		case ok1 && ok2 && val1 != val2:
			// Значения не карты, не равны.
			// Здесь код выглядит сложно, но лучше я придумать не смог, с учетом придуманной мной структуры хранения и необходимости введения в дальнейшем нескольких форматтеров)
			// 1. создаем срез, в котором будем хранить два объекта: первый с ключем и его значением до изменения, другой - после изменения
			changedValue := make([]KeyCharacteristics, 0, 2)
			//2. описываем характеристики ключа до внесенных изменений
			newKeyChar.IsValueChanged = false
			newKeyChar.HasDiff = false
			newKeyChar.Key = key
			newKeyChar.Status = deleted
			newKeyChar.Value = val1
			//3. закидываем его в срез
			changedValue = append(changedValue, newKeyChar)
			//4. описываем характеристики ключа после изменений - тут достаточно поменять статус и значение ключа
			newKeyChar.Status = added
			newKeyChar.Value = val2
			//5. закидываем его в срез
			changedValue = append(changedValue, newKeyChar)
			//6. а здесь мы формируем объект, который в себе содержит наименование ключа, срез с его значениями до и после изменений
			newKeyChar.IsValueChanged = true
			newKeyChar.Status = "changed"
			newKeyChar.Value = changedValue
			//7. и его уже закидываем в дифф
			resultDiff[key] = newKeyChar
		case ok1 && !ok2:
			// ключ есть только в data1
			newKeyChar.Key = key
			newKeyChar.HasDiff = false
			newKeyChar.Status = deleted
			newKeyChar.Value = val1
			resultDiff[key] = newKeyChar
		case !ok1 && ok2:
			// ключ есть только в data2
			newKeyChar.Key = key
			newKeyChar.HasDiff = false
			newKeyChar.Status = added
			newKeyChar.Value = val2
			resultDiff[key] = newKeyChar
		}
	}
	return resultDiff
}
