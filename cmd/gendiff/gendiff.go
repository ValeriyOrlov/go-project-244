package gendiff

import (
	"fmt"
	"maps"
	"slices"
)

func Gendiff(data1, data2 map[string]any) map[string]string {
	const (
		equal      = " "
		rowOfData1 = "-"
		rowOfData2 = "+"
	)

	diff := make(map[string]string)
	data1Keys := slices.Sorted(maps.Keys(data1))
	data2Keys := slices.Sorted(maps.Keys(data2))
	allKeys := append(data1Keys, data2Keys...)
	uniqueKeys := slices.Compact(allKeys)
	slices.Sort(uniqueKeys)

	for _, key := range uniqueKeys {
		val1, ok1 := data1[key]
		val2, ok2 := data2[key]
		row1 := fmt.Sprintf("%s: %v", key, val1)
		row2 := fmt.Sprintf("%s: %v", key, val2)

		switch {
		case ok1 && ok2 && val1 == val2:
			// Значения совпадают
			diff[row1] = equal
		case ok1 && ok2 && val1 != val2:
			// Значения отличаются
			diff[row1] = rowOfData1
			diff[row2] = rowOfData2
		case ok1 && !ok2:
			// ключ есть только в data1
			diff[row1] = rowOfData1
		case !ok1 && ok2:
			// ключ есть только в data2
			diff[row2] = rowOfData2
		}
	}
	return diff
}
