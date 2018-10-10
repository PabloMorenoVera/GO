package lexer

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {

	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex := NewLexer(r, File)

	for index := 0; index < 13; index++ {
		token, err := lex.Lex()
		fmt.Println(token, err)
	}
}
