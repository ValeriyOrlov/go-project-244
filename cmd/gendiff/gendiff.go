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
	for key, value := range data1 {
		if slices.Contains(data2Keys, key) && data2[key] == value {
			row := fmt.Sprintf("%s: %v", key, value)
			diff[row] = equal
		} else if slices.Contains(data2Keys, key) && data2[key] != value {
			row1 := fmt.Sprintf("%s: %v", key, value)
			row2 := fmt.Sprintf("%s: %v", key, data2[key])
			diff[row1] = rowOfData1
			diff[row2] = rowOfData2
		} else {
			row := fmt.Sprintf("%s: %v", key, value)
			diff[row] = rowOfData1
		}
	}
	for key := range data2 {
		if !slices.Contains(data1Keys, key) {
			row := fmt.Sprintf("%s: %v", key, data2[key])
			diff[row] = rowOfData2
		}
	}
	return diff
}
