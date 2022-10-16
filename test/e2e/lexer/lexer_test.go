package lexer

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Compiler2022/compilers-1-Belstowe/pkg/librust"
)

func TestCorrectCode(t *testing.T) {
	examples := []string{
		"../../../examples/factorial.rs",
		"../../../examples/fibonacci.rs",
		"../../../examples/gcd.rs",
	}

	for _, example_path := range examples {
		example_data, err := os.Open(example_path)
		if err != nil {
			t.Fatalf("[%s] Couldn't open file, reason: %v", example_path, err)
		}
		buf := bytes.NewBufferString("")
		librust.Parse(bufio.NewReader(example_data), nil, buf, true, false, false, false)
		for _, line := range strings.Split(buf.String(), "\n") {
			if strings.Contains(line, "error") {
				t.Errorf("[%s] %s", example_path, line)
			}
		}
	}
}
