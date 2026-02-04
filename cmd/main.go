package main

import (
	"code/formatters"
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
			&cli.BoolFlag{
				Name:    "stylish",
				Aliases: []string{"s"},
				Usage:   "stylish format",
				Value:   true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 2 {
				fmt.Println("file paths are not specified (use the -h flag for reference)")
				return nil
			}
			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)

			result := formatters.Formatters(path1, path2, cmd.Bool("stylish"))
			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
