package main

import (
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
			fmt.Println("test start")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
