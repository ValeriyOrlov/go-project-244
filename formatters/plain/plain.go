package formatters

import (
	"code/cmd/gendiff"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

func IsBool(v interface{}) bool {
	if v == nil {
		return false
	}
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

func Plain(data map[string]gendiff.KeyCharacteristics, path []string) string {
	result := []string{}
	keys := slices.Sorted(maps.Keys(data))
	keyStatuses := map[string]string{"added": " was added with value: ", "deleted": " was removed", "changed": " was updated. "}
	var removedVal interface{}
	var addedVal interface{}
	complexVal := "[complex value]"
	propertyStr := "Property "
	for _, key := range keys {
		if data[key].HasDiff {
			path = append(path, data[key].Key, ".")
			result = append(result, Plain(data[key].Value.(map[string]gendiff.KeyCharacteristics), path))
			if len(path) != 0 {
				path = path[0 : len(path)-2]
			} else {
				path = []string{}
			}
		} else {
			switch data[key].Status {
			case "added":
				if gendiff.IsMap(data[key].Value) {
					addedVal = complexVal
				} else if IsBool(data[key].Value) || data[key].Value == nil {
					addedVal = fmt.Sprintf(`%v`, data[key].Value)
				} else {
					addedVal = fmt.Sprintf(`'%v'`, data[key].Value)
				}
				result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), data[key].Key), keyStatuses[data[key].Status], fmt.Sprint(addedVal), "\n")
			case "deleted":
				result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), data[key].Key), keyStatuses[data[key].Status], "\n")
			case "changed":
				if gendiff.IsMap(data[key].Value.([]gendiff.KeyCharacteristics)[0].Value) {
					removedVal = complexVal
				} else if IsBool(data[key].Value.([]gendiff.KeyCharacteristics)[0].Value) || data[key].Value.([]gendiff.KeyCharacteristics)[0].Value == nil {
					removedVal = fmt.Sprintf(`%v`, data[key].Value.([]gendiff.KeyCharacteristics)[0].Value)
				} else {
					removedVal = fmt.Sprintf(`'%v'`, data[key].Value.([]gendiff.KeyCharacteristics)[0].Value)
				}
				if gendiff.IsMap(data[key].Value.([]gendiff.KeyCharacteristics)[1].Value) {
					addedVal = complexVal
				} else if IsBool(data[key].Value.([]gendiff.KeyCharacteristics)[1].Value) || data[key].Value.([]gendiff.KeyCharacteristics)[1].Value == nil {
					addedVal = fmt.Sprintf(`%v`, data[key].Value.([]gendiff.KeyCharacteristics)[1].Value)
				} else {
					addedVal = fmt.Sprintf(`'%v'`, data[key].Value.([]gendiff.KeyCharacteristics)[1].Value)
				}
				changedValue := fmt.Sprintf("From %v to %v.", removedVal, addedVal)
				result = append(result, propertyStr, fmt.Sprintf(`'%s%s'`, strings.Join(path, ""), data[key].Key), keyStatuses[data[key].Status], changedValue, "\n")
			}
		}
	}
	return strings.Join(result, "")
}
