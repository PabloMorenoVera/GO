package lexer

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {

	f, err := os.Open("prueba.txt")
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex := NewLexer(r, File)

	token, err := lex.Lex()
	fmt.Println(token.lexema, err)
}
