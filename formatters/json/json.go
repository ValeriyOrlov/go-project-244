package formatters

import (
	"code/cmd/gendiff"
	"encoding/json"
	"log"
)

func Json(data []gendiff.KeyCharacteristics) string {
	result, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		return "[]" // возвращаем пустой массив, чтобы не ломать формат
	}
	return string(result)
}
