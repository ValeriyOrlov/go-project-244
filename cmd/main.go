package main

import (
	"code/cmd/gendiff"
	"code/formatters"
	"code/parsers"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3" // imports as package "cli"
)

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

			data1, err1 := parsers.Parser(path1)
			if err1 != nil {
				log.Fatal(err1)
			}
			data2, err2 := parsers.Parser(path2)
			if err2 != nil {
				log.Fatal(err2)
			}
			diff := gendiff.Gendiff(data1, data2)
			result := formatters.Stylish(diff)
			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
