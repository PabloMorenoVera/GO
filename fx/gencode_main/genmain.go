package main

import (
	"bufio"
	"fmt"
	"fx/fxlex"
	"fx/gixgen"
	"fx/gixparse"
	"os"
)

func getTitle() string {
	return "\n\n# FX Lenguage Compiler. \n# Author: Pablo Moreno Vera\n"
}

func main() {
	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	fmt.Println(getTitle())

	r := bufio.NewReader(f)
	parser := gixparse.NewParser(fxlex.NewLexer(r, "lang.fx"))
	node, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error: %s, Parser: %s \n", err, node.PrintProg())
	}

	lstk := gixgen.NewStack()
	_ = lstk.ProgMain(node)
}
