package formatters

import (
	"code/cmd/gendiff"
	Plain "code/formatters/plain"
	Stylish "code/formatters/stylish"
	"code/parsers"
	"log"
)

func Formatters(path1 string, path2 string, stylish bool, plain bool) string {
	data1, err1 := parsers.Parser(path1)
	if err1 != nil {
		log.Fatal(err1)
	}
	data2, err2 := parsers.Parser(path2)
	if err2 != nil {
		log.Fatal(err2)
	}
	diff := gendiff.Gendiff(data1, data2)
	if plain {
		return Plain.Plain(diff, []string{})
	} else if stylish {
		return Stylish.Stylish(diff, 0)
	}
	return ""
}
