package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust"
	"github.com/urfave/cli/v2"
)

func main() {
	var input_path string
	var to_dump_tokens bool
	var to_dump_ast bool

	app := &cli.App{
		Name:  "ruster",
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
			&cli.BoolFlag{
				Name:        "dump-ast",
				Usage:       "Require parser to dump AST in stdout",
				Destination: &to_dump_ast,
			},
		},
		Action: func(c *cli.Context) error {
			var code io.Reader

			if !to_dump_ast && !to_dump_tokens {
				return nil
			}

			if input_path == "" {
				fmt.Println("Input Rust code (press Ctrl+D (Unix) or Ctrl+Z (Win) to interrupt):")
				code = bufio.NewReader(os.Stdin)
			} else {
				var err error
				code, err = os.Open(input_path)
				if err != nil {
					panic(err)
				}
			}

			if to_dump_tokens {
				librust.DumpTokens(code, os.Stdout)
			}

			if to_dump_ast {
				librust.Parse(code, os.Stdout)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
