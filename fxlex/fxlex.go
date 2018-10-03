package main

var file = "lang.fx"

type Token struct {
	lexema string
	tipo   string
	valor  string
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
}

func NewLexer(r RuneScanner, file string) (l *Lexer, err error) {
	l = &Lexer{line: 1}
	l.file = file
	l.r = r
	return l, nil
}

func (l *Lexer) get() (r rune, err error) {
	r, _, err = l.r.ReadRune()
	if err == nil {
		l.lastrune = r
		if r == '\n' {
			l.line++
		}
	}
	return r, err
}

func (l *Lexer) unget() (err error) {
	err = l.r.UnreadRune()
	if err == nil && l.lastrune == '\n' {
		l.line--
	}
	return err
}
