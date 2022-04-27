package parser

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust"
)

func TestCorrectCode(t *testing.T) {
	examples := []string{
		"../../../examples/find_substr.rs",
		"../../../examples/gcd.rs",
		"../../../examples/min.rs",
	}

	for _, example_path := range examples {
		example_data, err := os.Open(example_path)
		if err != nil {
			t.Fatalf("[%s] Couldn't open file, reason: %v", example_path, err)
		}
		buf := bytes.NewBufferString("")
		librust.Parse(bufio.NewReader(example_data), buf)
	}
}
