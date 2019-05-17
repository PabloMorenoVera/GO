package fxlex2

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	RuneEOF         = 0
	TokId   TokType = iota
	TokEOF
	TokFunc
	TokMain
	TokDecl
	TokIter
	TokIf
	TokVType
	TokTypeId

	// Punctuation Tokens
	TokLPar   TokType = '('
	TokRPar   TokType = ')'
	TokLCorch TokType = '['
	TokRCorch TokType = ']'
	TokComma  TokType = ','
	TokPC     TokType = ';'
	TokLKey   TokType = '{'
	TokRKey   TokType = '}'
	TokPunto  TokType = '.'

	// Int Operators Tokens
	TokSum  TokType = '+'
	TokRest TokType = '-'
	TokBar  TokType = '/'
	TokMin  TokType = '<'
	TokMax  TokType = '>'
	TokPorc TokType = '%'
	TokMul  TokType = '*'

	// Bool Operators Tokens
	TokOr  TokType = '|'
	TokAnd TokType = '&'
	TokNot TokType = '!'
	TokXOr TokType = '^'
)

type TokType rune

type Token struct {
	lexema string
	Type   TokType
	valor  int64
}

type RuneScanner interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

type Lexer struct {
	file     string
	line     int
	r        RuneScanner
	lastrune rune

	accepted []rune
	tokSaved *Token
}

func NewLexer(r RuneScanner, file string) (l *Lexer) {
	l = &Lexer{line: 1}
	l.file = file
	l.r = r
	return l
}

func (l *Lexer) get() (r rune) {
	var err error
	r, _, err = l.r.ReadRune()
	if err == nil {
		l.lastrune = r
	}
	if r == '\n' {
		l.line++
	}
	if err == io.EOF {
		l.lastrune = RuneEOF
		return RuneEOF
	}
	if err != nil {
		panic(err)
	}
	l.accepted = append(l.accepted, r)
	return r
}

func (l *Lexer) unget() {
	var err error
	if l.lastrune == RuneEOF {
		return
	}
	err = l.r.UnreadRune()
	if err == nil && l.lastrune == '\n' {
		l.line--
	}
	l.lastrune = unicode.ReplacementChar
	if len(l.accepted) != 0 {
		l.accepted = l.accepted[0 : len(l.accepted)-1]
	}
	if err != nil {
		panic(err)
	}
}

func (l *Lexer) accept() (tok string) {
	tok = string(l.accepted)
	if tok == "" && l.lastrune != RuneEOF {
		panic(errors.New("empty token"))
	}
	l.accepted = nil
	return tok
}

func (l *Lexer) readcomment() {
	for r := l.get(); ; r = l.get() {
		if r == '\n' {
			return
		}
	}
	return
}

func (l *Lexer) lexId() (t Token, err error) {
	r := l.get()
	if !unicode.IsLetter(r) {
		return t, errors.New("bad Id, should not happen")
	}
	isAlpha := func(ar rune) bool {
		return unicode.IsDigit(ar) || unicode.IsLetter(ar) || r == '_'
	}
	for r := l.get(); isAlpha(r); r = l.get() {
	}
	l.unget()
	switch string(l.accepted) {
	case "True", "False":
		t.Type = TokBool
	case "Coord":
		t.Type = TokCoord
	case "func":
		t.Type = TokFunc
	case "main":
		t.Type = TokMain
	case "iter":
		t.Type = TokIter
	case "if":
		t.Type = TokIf
	default:
		t.Type = TokId
	}
	t.lexema = l.accept()
	return t, nil
}

func (l *Lexer) lexNum() (t Token, err error) {
	const (
		Es    = "Ee"
		Signs = "+-"
	)
	hasDot := false
	r := l.get()
	if r == '.' {
		hasDot = true
		r = l.get()
	}
	for ; unicode.IsDigit(r); r = l.get() {
	}
	if r == '.' {
		if hasDot {
			return t, errors.New("bad float [" + l.accept() + "]")
		}
		hasDot = true
		for r = l.get(); unicode.IsDigit(r); r = l.get() {
		}
	}
	//fmt.Println(l.accepted)
	switch {
	case strings.ContainsRune(Es, r):
		r = l.get()
		if strings.ContainsRune(Signs, r) {
			r = l.get()
		}
	case hasDot:
		l.unget()
		break
	case !hasDot:
		if r == 'x' {
			for r = l.get(); ; r = l.get() {
				if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
					break
				}
			}
			l.unget()
			t.lexema = l.accept()
			t.valor, err = strconv.ParseInt(t.lexema, 0, 64)
			if err != nil {
				return t, errors.New("bad int [" + t.lexema + "]")
			}
			t.Type = TokInt
			return t, err
		} else {
			l.unget()
			t.lexema = l.accept()
			t.valor, err = strconv.ParseInt(t.lexema, 10, 64)
			if err != nil {
				return t, errors.New("bad int [" + t.lexema + "]")
			}
			t.Type = TokInt
			return t, nil
		}
	default:
		return t, errors.New("bad float [" + l.accept() + "]")
	}
	for r = l.get(); unicode.IsDigit(r); r = l.get() {
	}
	l.unget()
	t.lexema = l.accept()
	t.valor, err = strconv.ParseInt(t.lexema, 10, 64)
	if err != nil {
		return t, errors.New("bad int [" + t.lexema + "]")
	}
	t.Type = TokFloat
	return t, nil
}

func (l *Lexer) Peek() (t Token, err error) {
	t, err = l.Lex()
	if err == nil {
		l.tokSaved = &t
	}
	return t, nil
}

func (l *Lexer) Lex() (t Token, err error) {

	if l.tokSaved != nil {
		t = *l.tokSaved
		l.tokSaved = nil
		return t, nil
	}

	for r := l.get(); ; r = l.get() {

		if unicode.IsSpace(r) && r != '\n' || r == '\t' {
			l.accept()
			continue
		}
		switch r {
		case '+', '-', '/', '>', '<', '%', '*':
			switch r {
			case '/':
				r = l.get()
				if r == '/' {
					l.readcomment()
					l.accept()
					continue
				} else {
					l.unget()
					t.Type = TokOpInt
					t.lexema = l.accept()
					return t, err
				}
			case '*':
				r = l.get()
				if r != '*' {
					l.unget()
				}
			case '>', '<':
				r = l.get()
				if r != '=' {
					l.unget()
				}
			}
			t.Type = TokOpInt
			t.lexema = l.accept()
			return t, err
		case '(', ')', ',', '{', '}', ';', '[', ']', '.':
			t.Type = TokPunct
			t.lexema = l.accept()
			return t, nil
		case RuneEOF:
			t.Type = TokEOF
			l.accept()
			return t, err
		case '\n':
			l.accept()
			continue
		case ':', '=':
			if r == ':' && l.get() != '=' {
				l.unget()
				l.accept()
				err := fmt.Sprintf("bar rune %c %x", r, r)
				return t, errors.New(err)
			}
			t.Type = TokDecl
			t.lexema = l.accept()
			return t, err
		case '|', '&', '!', '^':
			t.lexema = l.accept()
			t.Type = TokOpBool
			return t, err
		}

		switch {
		case unicode.IsLetter(r):
			l.unget()
			t, err = l.lexId()
			return t, err
		case unicode.IsDigit(r):
			l.unget()
			t, err = l.lexNum()
			return t, err
		default:
			err := fmt.Sprintf("bad rune %c %x", r, r)
			return t, errors.New(err)
		}
	}
	return t, err
}
