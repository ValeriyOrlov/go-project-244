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
				Name:    "plain",
				Aliases: []string{"p"},
				Usage:   "shortcut for --format plain",
			},
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "shortcut for --format json",
			},
			&cli.BoolFlag{
				Name:    "stylish",
				Aliases: []string{"s"},
				Usage:   "shortcut for --format stylish",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() < 2 {
				return fmt.Errorf("missing file paths")
			}
			path1 := cmd.Args().Get(0)
			path2 := cmd.Args().Get(1)

			var format string
			switch {
			case cmd.Bool("plain"):
				format = "plain"
			case cmd.Bool("json"):
				format = "json"
			case cmd.Bool("stylish"):
				format = "stylish"
			default:
				format = "stylish"
			}

			result, err := formatters.Format(path1, path2, format)
			if err != nil {
				return fmt.Errorf("format error: %w", err)
			}
			fmt.Println(result)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
