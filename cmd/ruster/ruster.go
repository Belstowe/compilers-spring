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
	var output_path string
	var to_dump_tokens bool
	var to_dump_ast bool
	var to_dump_asm bool
	var verbose bool

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
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "ex.ll",
				Usage:       "Path for LLVM IR file",
				DefaultText: "ex.ll",
				Destination: &output_path,
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
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "Print info messages as well",
				Destination: &verbose,
			},
			&cli.BoolFlag{
				Name:        "dump-asm",
				Usage:       "Require parser to dump LLVM IR in stdout",
				Destination: &to_dump_asm,
			},
		},
		Action: func(c *cli.Context) error {
			var code io.Reader

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

			if output_path == "" {
				librust.Parse(code, nil, os.Stderr, to_dump_tokens, to_dump_ast, to_dump_asm, verbose)
			} else {
				wf, err := os.Create(output_path)
				if err != nil {
					panic(err)
				}
				defer wf.Close()
				librust.Parse(code, wf, os.Stderr, to_dump_tokens, to_dump_ast, to_dump_asm, verbose)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
