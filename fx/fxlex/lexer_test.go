package fxlex

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestLexer(t *testing.T) {
	var token Token

	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex := NewLexer(r, "lang.fx")

	for token.TokType != TokEOF {
		token, err = lex.Lex()
		lexema := token.lexema
		tokenType := convToken(token)
		value := token.valor
		fmt.Println(lexema, tokenType, value, err)
	}
	fmt.Println("----------------------- End lexer test -----------------------")
}

func TestToken(t *testing.T) {
	var token Token
	listToken := [16]TokType{TokId, TokInt, TokFloat, TokBool, TokCoord, TokOpInt, TokOpBool, TokEOF, TokEol, TokFunc, TokPunct, TokAsig, TokMain, TokIter, TokIf}

	for _, tok := range listToken {
		token.TokType = tok
		tokName := convToken(token)
		fmt.Println(tokName)
	}
	fmt.Println("----------------------- End Token test -----------------------")
}

func TestNewLexer(t *testing.T) {
	f, err := os.Open("lang.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex := NewLexer(r, "lang.fx")

	fmt.Println(lex)
	fmt.Println("----------------------- End NewLexer test -----------------------")
}

func TestLexid(t *testing.T) {
	var token Token

	f, err := os.Open("langId.fx")
	defer f.Close()
	if err != nil {
		print(err)
	}

	r := bufio.NewReader(f)
	lex := NewLexer(r, "langId.fx")

	for index := 0; index < 13; index++ {
		token, err = lex.lexID()
		fmt.Println(token)
	}
	fmt.Println("----------------------- End lexID test -----------------------")
}
