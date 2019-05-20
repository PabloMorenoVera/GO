package fxparser

import (
	"bufio"
	"fmt"
	"fx/fxlex2"
	"os"
	"testing"
)

func TestNewLexer(t *testing.T) {
	var token fxlex2.Token

	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	parser := NewParser(fxlex2.NewLexer(r, "lang.fx"))
	fmt.Println("Parser:", parser)

	for token.Type != fxlex2.TokEOF {
		token, err = parser.l.Lex()
		fmt.Println("Lexema:", token.Lexema, fmt.Sprintf("Type: %v", token.Type), "Valor:", token.Valor)
	}

}
