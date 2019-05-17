package fxparser

import (
	"errors"
	"fmt"
	"fx/fxlex2"
)

type Parser struct {
	l     *fxlex2.Lexer
	depth int
}

func NewParser(l *fxlex2.Lexer) *Parser {
	return &Parser{l, 0}
}

func (p *Parser) match(tT fxlex2.TokType) (t fxlex2.Token, e error, isMatch bool) {
	t, err := p.l.Peek()
	if err != nil {
		return fxlex2.Token{}, err, false
	}
	if t.Type != tT {
		return t, nil, false
	}
	t, err = p.l.Lex()
	return t, nil, true
}

// <PROGRAM> ::= <FUNCTION> <PROGRAM>   |  'eof'
func (p *Parser) Program() error {
	_, err, isEOF := p.match(fxlex2.TokEOF)
	if err != nil || !isEOF {
		return err
	}
	if err := p.Function(); err != nil {
		return err
	}
	return p.Program()
}

// <FUNCTION> ::= <HEADER>  '{' <STATEMENTS> '}'
func (p *Parser) Function() error {
	if err := p.Header(); err != nil {
		return err
	}
	_, err, isLCorch := p.match(fxlex2.TokLCorch)
	if err != nil || !isLCorch {
		return err
	}
	if err := p.Statement(); err != nil {
		return err
	}
	_, err, isRCorch := p.match(fxlex2.TokRCorch)
	if err != nil || !isRCorch {
		return err
	}
	return errors.New("bad Function")
}

// <HEADER> ::= 'func'  Id  '(' <OPT_PARAMS> ')'
func (p *Parser) Header() error {
	_, err, isfunc := p.match(fxlex2.TokFunc)
	if err != nil || !isfunc {
		return err
	}
	_, err, isID := p.match(fxlex2.TokId)
	if err != nil || !isID {
		return err
	}
	_, err, isLPar := p.match(fxlex2.TokLPar)
	if err != nil || !isLPar {
		return err
	}
	if err := p.OptParams(); err != nil {
		return err
	}
	_, err, isRPar := p.match(fxlex2.TokRPar)
	if err != nil || !isRPar {
		return err
	}
	return errors.New("bad Function")
}

// <OPT_PARAMS>  ::= TypeId  Id  <PARAMS>  |   Empty
func (p *Parser) OptParams() error {
	_, err, isTID := p.match(fxlex2.TokTypeId)
	if err != nil || !isTID { //Empty
		return err
	}
	_, err, isID := p.match(fxlex2.TokId)
	if err != nil || !isID {
		return err
	}
	return p.Params()
}

// <PARAMS>  ::= ','  TypeId  Id   <PARAMS>  |  Empty
func (p *Parser) Params() error {
	_, err, isComma := p.match(fxlex2.TokComma)
	if err != nil || !isComma { //Empty
		return err
	}
	_, err, isTID := p.match(fxlex2.TokTypeId)
	if err != nil || !isTID {
		return err
	}
	_, err, isID := p.match(fxlex2.TokId)
	if err != nil || !isID {
		return err
	}
	return p.Params()
}

// <STATEMENTS>    ::= 'iter' '(' Id  ':='  <EXPR> ';' <EXPR> ',' <EXPR> ')' '{' <STATEMENT> '}'  <STATEMENT>  |     Id   '(' <FUNC_CALL> ';'   <STATEMENT>    |     Empty
func (p *Parser) Statement() error {
	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokIter:
		_, err, isLPar := p.match(fxlex2.TokLPar)
		if err != nil || !isLPar {
			return err
		}
		_, err, isID := p.match(fxlex2.TokId)
		if err != nil || !isID {
			return err
		}
		_, err, isAsig := p.match(fxlex2.TokAsig)
		if err != nil || !isAsig {
			return err
		}
		if err := p.Expr(); err != nil {
			return err
		}
		_, err, isPC := p.match(fxlex2.TokPC)
		if err != nil || !isPC {
			return err
		}
		if err := p.Expr(); err != nil {
			return err
		}
		_, err, isComma := p.match(fxlex2.TokComma)
		if err != nil || !isComma {
			return err
		}
		if err := p.Expr(); err != nil {
			return err
		}
		_, err, isID2 := p.match(fxlex2.TokId)
		if err != nil || !isID2 {
			return err
		}
		_, err, isRPar := p.match(fxlex2.TokRPar)
		if err != nil || !isRPar {
			return err
		}
		_, err, isLCorch := p.match(fxlex2.TokLCorch)
		if err != nil || !isLCorch {
			return err
		}
		if err := p.Statement(); err != nil {
			return err
		}
		_, err, isRCorch := p.match(fxlex2.TokRCorch)
		if err != nil || !isRCorch {
			return err
		}
		return p.Statement()
	case fxlex2.TokId:
		_, err, isLPar := p.match(fxlex2.TokLPar)
		if err != nil || !isLPar {
			return err
		}
		if err := p.FuncCall(); err != nil {
			return err
		}
		_, err, isPC := p.match(fxlex2.TokPC)
		if err != nil || !isPC {
			return err
		}
		return p.Statement()
	default:
		return nil //Empty
	}
}

// <FUNC_CALL>  ::= <EXPR>  <ARGS>  ')'   |    ')'
func (p *Parser) FuncCall() error {
	_, err, isRPar := p.match(fxlex2.TokRPar)
	if err == nil && !isRPar {
		if err := p.Expr(); err != nil {
			return err
		}
		if err := p.Args(); err != nil {
			return err
		}
		_, err, isRPar := p.match(fxlex2.TokRPar)
		if err != nil || !isRPar {
			return err
		}
	} else if err != nil || !isRPar {
		return err
	}
	return nil
}

// <ARGS>  ::= ','  <EXPR>  <ARGS>  |   Empty
func (p *Parser) Args() error {
	_, err, isComma := p.match(fxlex2.TokTypeId)
	if err != nil || !isComma { //Empty
		return err
	}
	if err := p.Expr(); err != nil {
		return err
	}
	return p.Params()
}

// <EXPR> ::= int_literal  |  bool_literal  |  Id
func (p *Parser) Expr() error {
	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokInt, fxlex2.TokBool, fxlex2.TokId:
		t, err := p.l.Lex()
		fmt.Print(t)
		return err
	default:
		return errors.New("bad Expr")
	}
}
