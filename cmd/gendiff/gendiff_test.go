package gendiff_test

import (
	"code/cmd/gendiff"
	"reflect"
	"testing"
)

func TestIsMap(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want bool
	}{
		{"nil", nil, false},
		{"int", 42, false},
		{"string", "str", false},
		{"slice", []int{1, 2}, false},
		{"map[string]int", map[string]int{"a": 1}, true},
		{"map[int]string", map[int]string{1: "a"}, true},
		{"map[string]interface{}", map[string]interface{}{"a": 1}, true},
		{"struct", struct{}{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gendiff.IsMap(tt.v); got != tt.want {
				t.Errorf("IsMap(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestGendiff(t *testing.T) {
	// Вспомогательные функции для создания ожидаемых структур
	eq := func(name string, val interface{}) gendiff.KeyCharacteristics {
		return gendiff.KeyCharacteristics{Name: name, Value: val, Status: "equal"}
	}
	del := func(name string, val interface{}) gendiff.KeyCharacteristics {
		return gendiff.KeyCharacteristics{Name: name, Value: val, Status: "deleted"}
	}
	add := func(name string, val interface{}) gendiff.KeyCharacteristics {
		return gendiff.KeyCharacteristics{Name: name, Value: val, Status: "added"}
	}
	chg := func(name string, children []gendiff.KeyCharacteristics) gendiff.KeyCharacteristics {
		return gendiff.KeyCharacteristics{Name: name, Value: children, Status: "changed"}
	}

	tests := []struct {
		name   string
		data1  map[string]interface{}
		data2  map[string]interface{}
		expect []gendiff.KeyCharacteristics
	}{
		{
			name:   "both empty",
			data1:  map[string]interface{}{},
			data2:  map[string]interface{}{},
			expect: []gendiff.KeyCharacteristics{},
		},
		{
			name:   "first empty, second has key",
			data1:  map[string]interface{}{},
			data2:  map[string]interface{}{"a": 1},
			expect: []gendiff.KeyCharacteristics{add("a", 1)},
		},
		{
			name:   "first has key, second empty",
			data1:  map[string]interface{}{"a": 1},
			data2:  map[string]interface{}{},
			expect: []gendiff.KeyCharacteristics{del("a", 1)},
		},
		{
			name:  "identical keys and values",
			data1: map[string]interface{}{"a": 1, "b": "str", "c": true},
			data2: map[string]interface{}{"a": 1, "b": "str", "c": true},
			expect: []gendiff.KeyCharacteristics{
				eq("a", 1),
				eq("b", "str"),
				eq("c", true),
			},
		},
		{
			name:  "added and deleted keys",
			data1: map[string]interface{}{"a": 1, "b": 2},
			data2: map[string]interface{}{"b": 2, "c": 3},
			expect: []gendiff.KeyCharacteristics{
				del("a", 1),
				eq("b", 2),
				add("c", 3),
			},
		},
		{
			name:  "changed simple value",
			data1: map[string]interface{}{"a": 1},
			data2: map[string]interface{}{"a": 2},
			expect: []gendiff.KeyCharacteristics{
				del("a", 1),
				add("a", 2),
			},
		},
		{
			name:  "mixed simple changes and equal",
			data1: map[string]interface{}{"a": 1, "b": 2, "d": 4},
			data2: map[string]interface{}{"a": 1, "b": 3, "c": 5},
			expect: []gendiff.KeyCharacteristics{
				eq("a", 1),
				del("b", 2),
				add("b", 3),
				add("c", 5),
				del("d", 4),
			},
		},
		{
			name: "nested maps – both maps, equal",
			data1: map[string]interface{}{
				"nested": map[string]interface{}{"x": 10, "y": 20},
			},
			data2: map[string]interface{}{
				"nested": map[string]interface{}{"x": 10, "y": 20},
			},
			expect: []gendiff.KeyCharacteristics{
				chg("nested", []gendiff.KeyCharacteristics{
					eq("x", 10),
					eq("y", 20),
				}),
			},
		},
		{
			name: "nested maps – both maps, inner changed",
			data1: map[string]interface{}{
				"nested": map[string]interface{}{"x": 10, "y": 20},
			},
			data2: map[string]interface{}{
				"nested": map[string]interface{}{"x": 99, "y": 20},
			},
			expect: []gendiff.KeyCharacteristics{
				chg("nested", []gendiff.KeyCharacteristics{
					del("x", 10),
					add("x", 99),
					eq("y", 20),
				}),
			},
		},
		{
			name: "nested maps – one map, other not a map",
			data1: map[string]interface{}{
				"nested": map[string]interface{}{"a": 1},
			},
			data2: map[string]interface{}{
				"nested": 42,
			},
			expect: []gendiff.KeyCharacteristics{
				del("nested", map[string]interface{}{"a": 1}),
				add("nested", 42),
			},
		},
		{
			name: "deep nested mixed",
			data1: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"c": 3,
					"d": map[string]interface{}{
						"e": 5,
					},
				},
			},
			data2: map[string]interface{}{
				"a": 1,
				"b": map[string]interface{}{
					"c": 3,
					"d": map[string]interface{}{
						"e": 6,
						"f": 7,
					},
					"g": 8,
				},
				"h": 9,
			},
			expect: []gendiff.KeyCharacteristics{
				eq("a", 1),
				chg("b", []gendiff.KeyCharacteristics{
					eq("c", 3),
					chg("d", []gendiff.KeyCharacteristics{
						del("e", 5),
						add("e", 6),
						add("f", 7),
					}),
					add("g", 8),
				}),
				add("h", 9),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gendiff.Gendiff(tt.data1, tt.data2)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("Gendiff() = %v, want %v", got, tt.expect)
			}
		})
	}
}
