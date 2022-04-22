package input

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func ReadStdinUntilInterrupt() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString(0)
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}
	return text
}

func ReadFile(path string) string {
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(dat)
}
