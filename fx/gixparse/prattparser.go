package gixparse

import (
	"errors"
	"fmt"
	"fx/fxlex"
)

var precTab = map[rune]int{
	')':                  1,
	'|':                  10,
	'^':                  20,
	'&':                  30,
	'>':                  40,
	'<':                  40,
	rune(fxlex.TokOpInt): 40,
	'+':                  50,
	'-':                  50,
	'/':                  60,
	'*':                  60,
	'!':                  70,
	'(':                  80,
}

const defRbp = 0

func bindPow(tok fxlex.Token) int {
	if rbp, ok := precTab[rune(tok.Type)]; ok {
		return rbp
	}
	return defRbp
}

var leftTab = map[rune]bool{
	'^': true,
}

var unaryTab = map[rune]bool{
	'+': true,
	'-': true,
	'(': true,
	'*': true,
	'/': true,
	'<': true,
	'>': true,
	'|': true,
}

func (p *Parser) Nud(tok fxlex.Token) (expr *Expr, err error) {
	var rExpr *Expr
	var rbp int

	fmt.Sprintf("Nud: %d, %v \n", rbp, tok)
	if tok.Type == fxlex.TokLPar {
		expr, err = p.Expr(rbp)
		if err != nil {
			return nil, err
		}
		if _, err, isClosed := p.match(fxlex.TokRPar); err != nil {
			return nil, err
		} else if !isClosed {
			return nil, errors.New("unmatched parenthesis")
		}
		return expr, nil
	}

	expr = NewExpr(tok)
	rbp = bindPow(tok)
	rTok := rune(tok.Type)

	if rbp != defRbp {
		if !unaryTab[rTok] {
			errs := fmt.Sprintf("%s is not unary", tok.Type)
			return nil, errors.New(errs)
		}
		rExpr, err = p.Expr(rbp)
		if rExpr == nil {
			return nil, errors.New("unary oper. without operand")
		}
		expr.ERight = rExpr
	}
	return expr, nil
}

func (p *Parser) Led(left *Expr, tok fxlex.Token) (expr *Expr, err error) {
	var rbp int

	expr = NewExpr(tok)
	expr.ELeft = left
	rbp = bindPow(tok)

	if isleft := leftTab[rune(tok.Type)]; isleft {
		rbp -= 1
	}
	fmt.Sprintf("Led: %d, {{%v}} %v \n", rbp, left, tok)
	rExpr, err := p.Expr(rbp)
	if err != nil {
		return nil, err
	}
	if rExpr == nil {
		errs := fmt.Sprintf("missing operand for %s \n", tok.Type)
		return nil, errors.New(errs)
	}
	expr.ERight = rExpr
	return expr, nil
}

func (p *Parser) Expr(rbp int) (expr *Expr, err error) {
	var left *Expr

	s := fmt.Sprintf("Expr: %d", rbp)
	p.pushTrace(s)
	defer p.popTrace()

	tok, err := p.l.Peek()
	if err != nil {
		return expr, err
	}
	fmt.Sprintf("expr: Nud Lex: %v", tok)
	if tok.Type == fxlex.TokEOF {
		return expr, nil
	}
	p.l.Lex()
	if left, err = p.Nud(tok); err != nil {
		return nil, err
	}

	expr = left

	for {
		tok, err := p.l.Peek()
		if err != nil {
			return nil, err
		}
		if tok.Type == fxlex.TokEOF || tok.Type == fxlex.TokRPar {
			return expr, nil
		}
		if bindPow(tok) <= rbp {
			fmt.Sprintf("Not enough binding:")
			fmt.Sprintf("%d <= %d, %v \n", bindPow(tok), rbp, tok)
			return left, nil
		}
		p.l.Lex()
		fmt.Sprintf("expr: Led Lex: %v", tok)
		if left, err = p.Led(left, tok); err != nil {
			return nil, err
		}
		expr = left
	}
	return expr, err
}
