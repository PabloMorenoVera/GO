package gixgen

import (
	"bufio"
	"fmt"
	"fx/fxlex"
	"fx/gixparse"
	"os"
	"testing"
)

func TestCode(t *testing.T) {
	//var stk *LabelStack

	//stk = LabelStack{}
	//stk.NewStack()

	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}
	//fmt.Println("---------------------------------- Starting Test----------------------------------")
	r := bufio.NewReader(f)
	parser := gixparse.NewParser(fxlex.NewLexer(r, "lang.fx"))
	node, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error: %s, Parser: %s \n", err, node.PrintProg())
	}
	//else {
	//	fmt.Printf("Parser: %s \n", node.PrintProg())
	//}

	//fmt.Println("---------------------------------- End Test----------------------------------")
	//fmt.Println("---------------------------------- Starting Code Generation---------------------------------")

	lstk := NewStack()
	_ = lstk.ProgMain(node)
	//fmt.Println("Stack:", stk)
	//fmt.Println("---------------------------------- End Code Generation----------------------------------")
}
