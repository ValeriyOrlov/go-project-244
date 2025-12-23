package main

import (
	"code/formatters"
	"code/parsers"
	"context"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"

	"github.com/urfave/cli/v3" // imports as package "cli"
)

func gendiff(data1, data2 map[string]any) map[string]string {
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

func main() {
	cmd := &cli.Command{
		Name:  "gendiff",
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   `output format (default: "stylish")`,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("format") {
				fmt.Println("test format")
			}

			if cmd.Args().Len() < 2 {
				fmt.Println("file paths are not specified (use the -h flag for reference)")
				return nil
			}
			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)

			data1, _ := parsers.Parser(path1)
			data2, _ := parsers.Parser(path2)
			diff := gendiff(data1, data2)
			result := formatters.Stylish(diff)
			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
