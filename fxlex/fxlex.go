package lexer

import (
	"errors"
	"fmt"
	"io"
	"unicode"
)

const (
	File    = "lang.fx"
	RuneEOF = iota + 0x80
	TokId
	TokInt
	TokBool
	TokCoord
	TokOpInt
	TokOpBool
	TokEOF
	TokEol
	TokFunc
)

type tokType rune

type Token struct {
	lexema string
	tokType
	valor string
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
}

func NewLexer(r RuneScanner, file string) (l *Lexer) {
	l = &Lexer{line: 1}
	l.file = File
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

func (l *Lexer) Lex() (t Token, err error) {

	for r := l.get(); ; r = l.get() {
		if unicode.IsSpace(r) && r != '\n' {
			l.unget()
			t.lexema = l.accept()
			return t, err
		}
		switch r {
		case RuneEOF:
			t.tokType = TokEOF
			l.accept()
			return t, err
		case '\n':
			t.tokType = TokEol
			l.accept()
			return t, err
		case '/':
			r = l.get()
			if r == '/' {
				l.readcomment()
			}
			return t, err
		default:
			err := fmt.Sprintf("bad rune %c %x", r, r)
			return t, errors.New(err)
		}
	}
	return t, err
}
