package gixparse

import (
	"bufio"
	"fmt"
	"fx/fxlex"
	"os"
	"testing"
)

func TestLang(t *testing.T) {
	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}
	fmt.Println("---------------------------------- Starting Test----------------------------------")
	r := bufio.NewReader(f)
	parser := NewParser(fxlex.NewLexer(r, "lang.fx"))
	node, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error: %s, Parser: %s \n", err, node.PrintProg())
	} else {
		fmt.Printf("Parser: %s \n", node.PrintProg())
	}

	fmt.Println("---------------------------------- End Test----------------------------------")
}
