package fxparser

import (
	"bufio"
	"fx/fxlex2"
	"os"
	"testing"
)

func TestNewLexer(t *testing.T) {
	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)

	parser := NewParser(fxlex2.NewLexer(r, "lang.fx"))
}
