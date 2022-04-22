package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/input"
	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust"
	"github.com/urfave/cli/v2"
)

func main() {
	var input_path string
	var to_dump_tokens bool

	app := &cli.App{
		Name:  "rust-parser",
		Usage: "A simple Rust compiler using ANTLR",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "input",
				Aliases:     []string{"i"},
				Value:       "",
				Usage:       "Path to Rust code file for parsing",
				DefaultText: "read from terminal",
				Destination: &input_path,
			},
			&cli.BoolFlag{
				Name:        "dump-tokens",
				Usage:       "Require lexer to dump tokens in stdout",
				Destination: &to_dump_tokens,
			},
		},
		Action: func(c *cli.Context) error {
			var code string

			if input_path == "" {
				fmt.Println("Input Rust code (press Ctrl+D (Unix) or Ctrl+Z (Win) to interrupt):")
				code = input.ReadStdinUntilInterrupt()
			} else {
				code = input.ReadFile(input_path)
			}

			if to_dump_tokens {
				librust.DumpTokens(code, os.Stdout)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
