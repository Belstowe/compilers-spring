package symtab

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
	}

	for _, example_path := range examples {
		example_data, err := os.Open(example_path)
		if err != nil {
			t.Fatalf("[%s] Couldn't open file, reason: %v", example_path, err)
		}
		buf := bytes.NewBufferString("")
		librust.Parse(bufio.NewReader(example_data), buf, false, false, true)
		for _, line := range strings.Split(buf.String(), "\n") {
			if strings.Contains(line, "WARN") || strings.Contains(line, "ERROR") {
				t.Errorf("[%s] %s", example_path, line)
			}
		}
	}
}

func TestFlawedCode(t *testing.T) {
	examples := []string{
		"../../../examples/invalid/semantics/redefine_fn.rs",
		"../../../examples/invalid/semantics/redefine_vars.rs",
		"../../../examples/invalid/semantics/undefined.rs",
	}

	for _, example_path := range examples {
		example_data, err := os.Open(example_path)
		if err != nil {
			t.Fatalf("[%s] Couldn't open file, reason: %v", example_path, err)
		}
		buf := bytes.NewBufferString("")
		librust.Parse(bufio.NewReader(example_data), buf, false, false, false)

		if !strings.Contains(buf.String(), "ERROR") {
			t.Errorf("[%s] Should have delivered ERROR message, but hasn't.", example_path)
		}
	}
}

func TestSlightlyFlawedCode(t *testing.T) {
	examples := []string{
		"../../../examples/invalid/semantics/redefine_std.rs",
		"../../../examples/invalid/semantics/redefine_upperscope.rs",
	}

	for _, example_path := range examples {
		example_data, err := os.Open(example_path)
		if err != nil {
			t.Fatalf("[%s] Couldn't open file, reason: %v", example_path, err)
		}
		buf := bytes.NewBufferString("")
		librust.Parse(bufio.NewReader(example_data), buf, false, false, false)

		if !strings.Contains(buf.String(), "WARN") {
			t.Errorf("[%s] Should have delivered WARN message, but hasn't.", example_path)
		}
	}
}
