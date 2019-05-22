package fxparser

import (
	"errors"
	"fmt"
	"fx/fxlex2"
	"os"
	"strings"
)

type Parser struct {
	l     *fxlex2.Lexer
	depth int
}

func NewParser(l *fxlex2.Lexer) *Parser {
	return &Parser{l, 0}
}

func (p *Parser) pushTrace(tag string) {
	DebugDesc := true
	if DebugDesc {
		tabs := strings.Repeat("\t", p.depth)
		fmt.Fprintf(os.Stderr, "%s %s \n", tabs, tag)
	}
	p.depth++
}

func (p *Parser) popTrace() {
	p.depth--
}

func (p *Parser) Parse() error {
	p.pushTrace("Parse")
	defer p.popTrace()
	if err := p.Program(); err != nil {
		return err
	}
	return nil
}

func (p *Parser) match(tT fxlex2.TokType) (t fxlex2.Token, e error, isMatch bool) {
	t, err := p.l.Peek()
	fmt.Println("t:", t.Type, "tT:", tT)
	if err != nil {
		return fxlex2.Token{}, err, false
	}
	if t.Type != tT {
		return t, nil, false
	}
	t, err = p.l.Lex()
	return t, nil, true
}

// <PROGRAM> ::= <FUNCTION> <PROGRAM>   |
//								'type'	'record'	'Id'	'('	<EXPR>	<EXPR>	','	<EXPR>	<EXPR>	<DECL_PARAMS>						|
//								'eof'
func (p *Parser) Program() error {
	p.pushTrace("Program")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokEOF:
		_, err = p.l.Lex()
		return err
	case fxlex2.TokTypeId:
		_, err = p.l.Lex()
		_, err, isRecord := p.match(fxlex2.TokRecord)
		if err != nil || !isRecord {
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
		if err := p.Expr(); err != nil {
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
		if err := p.Expr(); err != nil {
			return err
		}
		if err := p.DeclParams(); err != nil {
			return err
		}
		return p.Program()
	default:
		if err := p.Function(); err != nil {
			return err
		}
		return p.Program()
	}
}

//	<DECL_PARAMS>		::=	')'	|
//											','	<EXPR>	<EXPR>	')'
func (p *Parser) DeclParams() error {
	p.pushTrace("Func_Call")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokRPar:
		_, err = p.l.Lex()
		return err
	case fxlex2.TokComma:
		_, err = p.l.Lex()
		if err := p.Expr(); err != nil {
			return err
		}
		if err := p.Expr(); err != nil {
			return err
		}
		_, err, isRPar := p.match(fxlex2.TokRPar)
		if err != nil || !isRPar {
			return err
		}
		return err
	default:
		return errors.New("bad DeclParams")
	}
}

// <FUNCTION> ::= <HEADER>  '{' <STATEMENTS> '}'
func (p *Parser) Function() error {
	p.pushTrace("Function")
	defer p.popTrace()

	if err := p.Header(); err != nil {
		return err
	}
	_, err, isLKey := p.match(fxlex2.TokLKey)
	if err != nil || !isLKey {
		return err
	}
	if err := p.Statement(); err != nil {
		return err
	}
	_, err, isRKey := p.match(fxlex2.TokRKey)
	if err != nil || !isRKey {
		return err
	}
	return errors.New("bad Function")
}

// <HEADER> ::= 'func'  Id  '(' <OPT_PARAMS> ')'
func (p *Parser) Header() error {
	p.pushTrace("Header")
	defer p.popTrace()

	_, err, isFUNC := p.match(fxlex2.TokFunc)
	if err != nil || !isFUNC {
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
	return err
}

// <OPT_PARAMS>  ::= TypeId  Id  <PARAMS>  |   Empty
func (p *Parser) OptParams() error {
	p.pushTrace("Opt_Params")
	defer p.popTrace()

	_, err, isTID := p.match(fxlex2.TokId) // Sería TokTypeId. Mirar
	if err != nil || !isTID {              //Empty
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
	p.pushTrace("Params")
	defer p.popTrace()

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

// <STATEMENTS>    ::= 'iter' '(' Id  ':='  <EXPR> ';' <EXPR> ',' <EXPR> ')' '{' <STATEMENT> '}'  <STATEMENT>  |
//												Id	<DECL>	|
//												Empty
func (p *Parser) Statement() error {
	p.pushTrace("Statement")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokIter:
		_, err = p.l.Lex()
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
		_, err, isRPar := p.match(fxlex2.TokRPar)
		if err != nil || !isRPar {
			return err
		}
		_, err, isLKey := p.match(fxlex2.TokLKey)
		if err != nil || !isLKey {
			return err
		}
		if err := p.Statement(); err != nil {
			return err
		}
		_, err, isRKey := p.match(fxlex2.TokRKey)
		if err != nil || !isRKey {
			return err
		}
		return p.Statement()
	case fxlex2.TokId:
		_, err = p.l.Lex()
		if err := p.Decl(); err != nil {
			return err
		}
		return p.Statement()
	default:
		return nil //Empty
	}
}

// <DECL>			::=		Id	','	|
//									'='	<EXPR>	';'	|
//									'(' <FUNC_CALL> ';'   <STATEMENTS>
func (p *Parser) Decl() error {
	p.pushTrace("Func_Call")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokId:
		_, err = p.l.Lex()
		_, err, isComma := p.match(fxlex2.TokComma)
		if err != nil || !isComma {
			return err
		}
	case fxlex2.TokEqual:
		_, err = p.l.Lex()
		if err := p.Expr(); err != nil {
			return err
		}
		_, err, isPC := p.match(fxlex2.TokPC)
		if err != nil || !isPC {
			return err
		}
	case fxlex2.TokLPar:
		_, err = p.l.Lex()
		if err := p.FuncCall(); err != nil {
			return err
		}
		_, err, isPC := p.match(fxlex2.TokPC)
		if err != nil || !isPC {
			return err
		}
		if err := p.Statement(); err != nil {
			return err
		}
		return err
	default:
		return errors.New("bad Decl")
	}
	return nil
}

// <FUNC_CALL>  ::= <EXPR>  <ARGS>  ')'   |    ')'
func (p *Parser) FuncCall() error {
	p.pushTrace("Func_Call")
	defer p.popTrace()

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
	p.pushTrace("Args")
	defer p.popTrace()

	_, err, isComma := p.match(fxlex2.TokTypeId)
	if err != nil || !isComma { //Empty
		return err
	}
	if err := p.Expr(); err != nil {
		return err
	}
	return p.Params()
}

// <EXPR> ::= int_literal  |  bool_literal  |  Id		|	Empty
func (p *Parser) Expr() error {
	p.pushTrace("Expr")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return err
	}
	switch t.Type {
	case fxlex2.TokInt:
		_, err := p.l.Lex()
		if err != nil {
			return err
		}
		t, err := p.l.Peek()
		if err != nil {
			return err
		}
		switch t.Type {
		case fxlex2.TokSum, fxlex2.TokRest, fxlex2.TokMul, fxlex2.TokBar, fxlex2.TokMax, fxlex2.TokMin, fxlex2.TokOpInt, fxlex2.TokPot, fxlex2.TokPorc:
			_, err := p.l.Lex()
			if err != nil {
				return err
			}
			return p.Expr()
		default:
			return err
		}
	case fxlex2.TokBool:
		_, err := p.l.Lex()
		if err != nil {
			return err
		}
		t, err := p.l.Peek()
		if err != nil {
			return err
		}

		switch t.Type {
		case fxlex2.TokOr, fxlex2.TokAnd, fxlex2.TokNot, fxlex2.TokXOr:
			_, err := p.l.Lex()
			if err != nil {
				return err
			}
			return p.Expr()
		default:
			return nil
		}
	case fxlex2.TokId:
		_, err := p.l.Lex()
		if err != nil {
			return err
		}
		t, err := p.l.Peek()
		if err != nil {
			return err
		}
		switch t.Type {
		case fxlex2.TokOr, fxlex2.TokAnd, fxlex2.TokNot, fxlex2.TokXOr, fxlex2.TokSum, fxlex2.TokRest, fxlex2.TokMul, fxlex2.TokBar, fxlex2.TokMax, fxlex2.TokMin, fxlex2.TokOpInt, fxlex2.TokPot, fxlex2.TokPorc:
			_, err := p.l.Lex()
			if err != nil {
				return err
			}
			return p.Expr()
		default:
			return nil
		}
	default:
		return nil
	}
}
