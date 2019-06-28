package gixparse

import (
	"errors"
	"fmt"
	"fx/fxlex"
	"fx/gixsymb"
	"math"
	"os"
	"strings"
)

type Parser struct {
	l      *fxlex.Lexer
	depth  int
	stkEnv gixsymb.StkEnv
}

type Const struct {
	cType fxlex.TokType
	val   float64
}

var sConst = map[string]Const{
	"Pi": Const{fxlex.TokFloat, math.Pi},
}

func (p *Parser) initSymbs() error {
	for k, v := range sConst {
		s, err := p.stkEnv.NewSymb(k, gixsymb.SConst)
		if err != nil {
			return err
		}
		s.TokKind = int(v.cType)
		s.FloatVal = v.val
		s.FName = "builtin"
	}
	return nil
}

func NewParser(l *fxlex.Lexer) *Parser {
	p := &Parser{l, 0, nil}
	p.stkEnv.PushEnv("Init")

	// Lenguage functions added to symbols table
	p.stkEnv.NewSymb("circle", gixsymb.SFunc)
	p.stkEnv.NewSymb("rect", gixsymb.SFunc)

	return p
}

func (p *Parser) pushTrace(tag string) {
	DebugDesc := false
	if DebugDesc {
		tabs := strings.Repeat("\t", p.depth)
		fmt.Fprintf(os.Stderr, "%s %s \n", tabs, tag)
	}
	p.depth++
}

func (p *Parser) popTrace() {
	p.depth--
}

func (p *Parser) recover(tok fxlex.TokType) error {
	token := fxlex.TokInit

	switch tok {
	case fxlex.TokLPar:
		token = fxlex.TokRPar
	case fxlex.TokLKey:
		token = fxlex.TokRKey
	default:
		return errors.New("Syntax Error: Not recovery")
	}

	t, err := p.l.Lex()
	if err != nil {
		return errors.New("Syntax Error: Not recovery")
	}
	for t.Type != token {
		t, err = p.l.Lex()
		if err != nil {
			return errors.New("Syntax Error: Not recovery")
		}
	}
	return errors.New("Syntax Error: Not found Parenthesis or Key")
}

func (p *Parser) match(tT fxlex.TokType) (t fxlex.Token, e error, isMatch bool) {
	t, err := p.l.Peek()
	if err != nil {
		return fxlex.Token{}, err, false
	}
	if t.Type != tT {
		return t, errors.New("bad TokenType Found"), false
	}
	t, err = p.l.Lex()
	return t, nil, true
}

func (p *Parser) Parse() (nprogram *Program, err error) {
	p.pushTrace("Parse")
	defer p.popTrace()

	nprogram, err = p.Program()
	if err != nil {
		return nprogram, err
	}

	return nprogram, nil
}

// <PROG> ::= <AUX>	'eof'
func (p *Parser) Program() (nprogram *Program, err error) {
	p.pushTrace("Program")
	defer p.popTrace()

	nprogram = NewProgram()
	nprogram.C, nprogram.depth = "Program", 0

	var parray []*Auxiliar
	nauxiliar, err := p.Auxiliar(parray, nprogram) // SENTENCE
	if err != nil {
		return nil, err
	}
	nprogram.Auxiliar = nauxiliar

	t, err, isEOF := p.match(fxlex.TokEOF)
	if err != nil || !isEOF {
		return nil, errors.New(fmt.Sprintf("Error: Bad %s, got %s", fxlex.TokEOF, t.Type))
	}

	return nprogram, nil
}

// <AUX>	::=	'func'	<FUNCTION>	<AUX>	|
//						'type'	<DECLARATION>	<AUX>	|
//						Empty
func (p *Parser) Auxiliar(parray []*Auxiliar, nprogram *Program) (array []*Auxiliar, err error) {
	p.pushTrace("Auxiliar")
	defer p.popTrace()

	nauxiliar := NewAuxiliar()
	nauxiliar.C, nauxiliar.depth = "Auxiliar", nprogram.depth+1

	t, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case fxlex.TokFunc: // Es FUNCTION
		_, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Missing TokFunc in Auxiliar")
		}

		var farray []*Function
		function, err := p.Function(farray, nprogram)
		if err != nil {
			return nil, err
		}
		nauxiliar.Function = function
		array = append(parray, nauxiliar)

		endarray, err := p.Auxiliar(array, nprogram)
		if err != nil {
			return nil, err
		}

		return endarray, nil
	case fxlex.TokTypeID: // Es DECLARATION
		t, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Missing TypeID in Auxiliar")
		}

		var darray []*Declaration
		declaration, err := p.Declaration(darray, nprogram) // DECLARATION
		if err != nil {
			return nil, err
		}
		nauxiliar.Declaration = declaration
		array = append(parray, nauxiliar)

		endarray, err := p.Auxiliar(array, nprogram)
		if err != nil {
			return nil, err
		}

		return endarray, nil
	default: // Empty
		return parray, nil
	}

}

// <DECLARATION>	::=	'record'	'Id'	'('	<PARAMS>	')'	|
//										Empty
func (p *Parser) Declaration(darray []*Declaration, nprogram *Program) (declaration []*Declaration, err error) {
	p.pushTrace("Declaration")
	defer p.popTrace()

	_, err, isRecord := p.match(fxlex.TokRecord) // Declaration
	if err != nil || !isRecord {                 // Es Empty
		return darray, nil
	}

	ndecl := NewDeclaration()
	ndecl.C = "Declaration"
	ndecl.depth = nprogram.depth + 1

	t, err, isID := p.match(fxlex.TokTypeID)
	if err != nil || !isID {
		return nil, errors.New("Error: Wrong TypeId in Declaration")
	}
	ndecl.Id = gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Var"}

	s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SVar)
	if err != nil || s == nil {
		return nil, err
	}

	_, err, isLPar := p.match(fxlex.TokLPar) // '('
	if err != nil || !isLPar {
		return nil, errors.New("Error: Missing Left Parenthesis in Declaration")
	}

	dparams, err := p.Params() // PARAMS
	if err != nil {
		return nil, err
	}
	ndecl.Params = dparams

	_, err, isRPar := p.match(fxlex.TokRPar) // ')'
	if err != nil || !isRPar {
		return nil, errors.New("Error: Missing Right Parenthesis in Declaration")
	}

	darray = append(darray, ndecl)
	return darray, nil
}

// <FUNCTION> ::=	Id  '(' <PARAMS> ')'  '{' <SENTENCE> '}'	<FUNCTION>	|
//								Empty
func (p *Parser) Function(farray []*Function, nprogram *Program) (function []*Function, err error) {
	p.pushTrace("Function")
	defer p.popTrace()

	t, err, isID := p.match(fxlex.TokId) // Func
	if err != nil || !isID {             // Es Empty
		return farray, nil
	}
	nfunc := NewFunction() // Create FUNCTION AST
	nfunc.C = "Function"
	nfunc.depth = nprogram.depth + 1
	nfunc.Id = t.Lexema

	// Create the Args and function body ambit
	s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SFunc)
	if err != nil || s == nil {
		return nil, err
	}
	p.stkEnv.PushEnv("Function")
	defer p.stkEnv.PopEnv("Function")

	_, err, isLPar := p.match(fxlex.TokLPar) // '('
	if err != nil || !isLPar {
		return nil, errors.New("Error: Missing Left Parenthesis in Function")
	}
	fparams, err := p.Params() // PARAMS
	if err != nil {
		return nil, err
	}
	nfunc.Params = fparams

	_, err, isRPar := p.match(fxlex.TokRPar) // ')'
	if err != nil || !isRPar {
		return nil, errors.New("Error: Missing Right Parenthesis in Function")
	}
	_, err, isLKey := p.match(fxlex.TokLKey) // '{'
	if err != nil || !isLKey {
		return nil, errors.New("Error: Missing Left Key in Function")
	}
	var sarray []*Sentence
	nsentence, err := p.Sentence(sarray, nfunc.depth) // SENTENCE
	if err != nil {
		return nil, err
	}
	nfunc.Sentence = nsentence

	_, err, isRKey := p.match(fxlex.TokRKey) // '}'
	if err != nil || !isRKey {
		return nil, errors.New("Error: Missing Right Key in Function")
	}

	farray = append(farray, nfunc)
	farray, err = p.Function(farray, nprogram) // FUNCTION
	if err != nil {
		return farray, err
	}
	return farray, nil
}

// <PARAMS>  ::= TypeId  Id  <OPT_PARAMS>  |   Empty
func (p *Parser) Params() (params []string, err error) {
	p.pushTrace("Opt_Params")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	//<OPT_PARAMS>  ::= TypeId  Id  <PARAMS>
	case fxlex.TokTypeID: // TypeId
		t, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Wrong TyppeID in Params")
		}
		params = append(params, t.Lexema)

		t, err, isID := p.match(fxlex.TokId) // Id
		if err != nil || !isID {
			return nil, errors.New("Error: Wrong ID in Params")
		}
		params = append(params, t.Lexema)
		s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SVar)
		if err != nil || s == nil {
			return nil, errors.New("Stack error: Param ID already declared")
		}

		params, err := p.OptParams(params) // OPT_PARAMS
		if params == nil {
			return nil, err
		}
		return params, err
	//<OPT_PARAMS>	::=	Empty
	default:
		return params, nil // Empty
	}
}

// <OPT_PARAMS>  ::= ','  TypeId  Id   <OPT_PARAMS>  |  Empty
func (p *Parser) OptParams(params []string) (optpar []string, err error) {
	p.pushTrace("Params")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	case fxlex.TokComma: // ','
		_, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Missing Comma in OptParams")
		}
		t, err, isTID := p.match(fxlex.TokTypeID) // TypeId
		if err != nil || !isTID {
			return nil, errors.New("Error: Wrong TypeID in OptParams")
		}
		params := append(params, t.Lexema)

		t, err, isID := p.match(fxlex.TokId) // Id
		if err != nil || !isID {
			return nil, errors.New("Error: Wrong ID in OptParams")
		}
		params = append(params, t.Lexema)
		s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SVar)
		if err != nil || s == nil {
			return nil, errors.New("Stack error: Param ID already declared")
		}

		optpar, err := p.OptParams(params) // OPT_PARAMS
		if err != nil {
			return params, err
		}
		return optpar, err
	default:
		return params, nil // Empty
	}
}

// <SENTENCE>    ::= 'iter' '(' Id  ':='  <ATOM> ';' <ATOM> ',' <ATOM> ')' '{' <SENTENCE> '}'	|
//										TypeId	Id	';'	|
//										Id	<ASIGFCALL>	|
//										Empty
func (p *Parser) Sentence(sarray []*Sentence, depth int) (sentence []*Sentence, err error) {
	p.pushTrace("Sentence")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return nil, err
	}
	switch t.Type {
	//<SENTENCE> ::= 'iter' '(' <INIT> <ATOM> ',' <ATOM> ')' '{' <SENTENCE> '}'  <SENTENCE>
	case fxlex.TokIter:
		nsent := NewSentence() // Create SENTENCE AST
		nsent.C = "Sentence"
		nsent.depth = depth + 1

		niter := NewIter() // Create ITER AST
		niter.C = "Iter"
		niter.depth = depth + 2
		nsent.Iter = niter
		sarray = append(sarray, nsent)

		t, err = p.l.Lex()
		niter.Id = t.Lexema

		_, err, isLPar := p.match(fxlex.TokLPar) // '('
		if err != nil || !isLPar {
			return nil, errors.New("Error: Missing Left Parenthesis in Iter")
		}

		p.stkEnv.PushEnv("Iter")

		nasig, err := p.Init(depth) // INIT
		if err != nil {
			return nil, err
		}
		niter.Init = nasig

		t, err := p.l.Lex() // ATOM
		if err != nil {
			return nil, err
		}
		niter.Atom1 = gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Var"}

		_, err, isComma := p.match(fxlex.TokComma) // ','
		if err != nil || !isComma {
			return nil, errors.New("Error: Missing Comma in Iter")
		}

		t, err = p.l.Lex() // ATOM
		if err != nil {
			return nil, err
		}
		niter.Atom2 = gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Var"}

		_, err, isRPar := p.match(fxlex.TokRPar) // ')'
		if err != nil || !isRPar {
			return nil, errors.New("Error: Missing Right Parenthesis in Iter")
		}
		_, err, isLKey := p.match(fxlex.TokLKey) // '{'
		if err != nil || !isLKey {
			return nil, errors.New("Error: Missing Left Key in Iter")
		}

		var iarray []*Sentence
		nsentence, err := p.Sentence(iarray, depth) // SENTENCE
		if err != nil {
			return nsentence, err
		}
		niter.Sentence = nsentence

		_, err, isRKey := p.match(fxlex.TokRKey) // '}'
		if err != nil || !isRKey {
			return nil, errors.New("Error: Missing Right Key in Iter")
		}

		p.stkEnv.PopEnv("Iter") // Pop iter environment

		sarray, err = p.Sentence(sarray, depth) // SENTENCE
		if err != nil {
			return nsentence, err
		}
		return sarray, nil
	//<SENTENCE>	::=	TypeId	Id	';'	<SENTENCE>
	case fxlex.TokTypeID:
		nsent := NewSentence() // Create SENTENCE AST
		nsent.C = "Sentence"
		nsent.depth = depth + 1

		t, err = p.l.Lex() // TypeId

		t, err, isID := p.match(fxlex.TokId) // Id
		if err != nil || !isID {
			return nil, errors.New("Error: Wrong ID in LocalDeclaration")
		}
		nsent.Decl = gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Var"}

		s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SVar)
		if err != nil || s == nil {
			return nil, errors.New("Stack error: Can't push new item ->" + t.Lexema)
		}

		_, err, isPC := p.match(fxlex.TokPC) // ';'
		if err != nil || !isPC {
			return nil, errors.New("Error: Missing Dot&Comma in LocalDeclaration")
		}

		sarray = append(sarray, nsent)

		sarray, err = p.Sentence(sarray, depth) // SENTENCE
		if err != nil {
			return sarray, err
		}
		return sarray, nil
	//<SENTENCE>	::=	Id	<ASIGFCALL>	|
	case fxlex.TokId:
		nsent := NewSentence() // Create SENTENCE AST
		nsent.C = "Sentence"
		nsent.depth = depth + 1

		t, err, isPC := p.match(fxlex.TokId) // Id
		if err != nil || !isPC {
			return nil, errors.New("Error: Missing Dot&Comma in LocalDeclaration")
		}

		// Get the variable name not the id and the fields
		id := strings.Split(t.Lexema, ".")
		var s *gixsymb.Symb
		if len(id) == 1 {
			s = p.stkEnv.GetSymb(t.Lexema)
		} else {
			s = p.stkEnv.GetSymb(id[0])
		}
		if s == nil {
			return nil, errors.New("Stack error: ID not declared -> " + t.Lexema)
		}

		nsent, err = p.AsigFcall(nsent, depth, t)
		if err != nil {
			return nil, err
		}
		sarray = append(sarray, nsent)

		sarray, err = p.Sentence(sarray, depth) // SENTENCE
		if err != nil {
			return sarray, err
		}

		return sarray, nil
	//<SENTENCE>	::=	TokIf	'('	<EXPR>	')'	'{'	<SENTENCE>	'}'	|
	case fxlex.TokIf:
		nsent := NewSentence() // Create SENTENCE AST
		nsent.C = "Sentence"
		nsent.depth = depth + 1

		nif := NewSif() // Create ITER AST
		nif.C = "If"
		nif.depth = depth + 2
		nsent.Sif = nif

		t, err = p.l.Lex()
		nif.Id = t.Lexema

		_, err, isLPar := p.match(fxlex.TokLPar) // '('
		if err != nil || !isLPar {
			return nil, errors.New("Error: Missing Left Parenthesis in If")
		}

		expr, err := p.Expr(1) // ATOM
		if err != nil {
			return nil, err
		}
		nif.Expr = expr

		_, err, isRPar := p.match(fxlex.TokRPar) // ')'
		if err != nil || !isRPar {
			return nil, errors.New("Error: Missing Right Parenthesis in If")
		}

		_, err, isLKey := p.match(fxlex.TokLKey) // '{'
		if err != nil || !isLKey {
			return nil, errors.New("Error: Missing Left Key in If")
		}

		var iarray []*Sentence
		nsentence, err := p.Sentence(iarray, depth) // SENTENCE
		if err != nil {
			return nsentence, err
		}
		nif.Sentence = nsentence

		_, err, isRKey := p.match(fxlex.TokRKey) // '}'
		if err != nil || !isRKey {
			return nil, errors.New("Error: Missing Right Parenthesys in If")
		}

		sarray = append(sarray, nsent)
		sarray, err = p.Sentence(sarray, depth) // SENTENCE
		if err != nil {
			return sarray, err
		}
		return sarray, nil
		//<SENTENCE>	::=	TokElse	'{'	<SENTENCE>	'}'
	case fxlex.TokElse:
		nsent := NewSentence() // Create SENTENCE AST
		nsent.C = "Sentence"
		nsent.depth = depth + 1

		nelse := NewSelse() // Create ITER AST
		nelse.C = "Else"
		nelse.depth = depth + 2
		nsent.Selse = nelse

		t, err = p.l.Lex()
		nelse.Id = t.Lexema

		_, err, isLKey := p.match(fxlex.TokLKey) // '{'
		if err != nil || !isLKey {
			return nil, errors.New("Error: Missing Left Key in Else")
		}

		var earray []*Sentence
		nsentence, err := p.Sentence(earray, depth) // SENTENCE
		if err != nil {
			return nsentence, err
		}
		nelse.Sentence = nsentence

		_, err, isRKey := p.match(fxlex.TokRKey) // '}'
		if err != nil || !isRKey {
			return nil, errors.New("Error: Missing Right Key in Else")
		}

		sarray = append(sarray, nsent)
		sarray, err = p.Sentence(sarray, depth) // SENTENCE
		if err != nil {
			return sarray, err
		}
		return sarray, nil
	default:
		return sarray, nil //Empty
	}
}

// <ASIGFCALL>	::=	'='	<EXPR>	';'	\
//									'('	<ATOM>	<ARGS>	')'	';'	\
//									Empty
func (p *Parser) AsigFcall(nsent *Sentence, depth int, r fxlex.Token) (sentence *Sentence, err error) {
	p.pushTrace("AsigFcall")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return nil, err
	}

	switch t.Type {
	//<ASIGFCALL> ::= '='	<EXPR>	';'
	case fxlex.TokEqual:
		nasig := NewAsig() // Create Asig AST
		nasig.C = "Asig"
		nasig.depth = depth + 1
		nasig.Id = gixsymb.Symb{r.Lexema, 0, int(r.Type), float64(r.Valor), "Var"}

		t, err = p.l.Lex() // '='
		if err != nil {
			return nil, errors.New("Error: Missing Equal in Asignation")
		}

		t, err := p.l.Peek()
		if err != nil {
			return nil, err
		}

		if t.Type == fxlex.TokLCorch { // Coords
			atom, err := p.Atom() // ATOM
			if err != nil {
				return nil, errors.New("Error: Missing Left Bracket in Asignation for Coord")
			}
			nasig.Coord = atom
		} else {
			expr, err := p.Expr(1) // EXPR

			if expr.Tok.Type == fxlex.TokId {
				s := p.stkEnv.GetSymb(expr.Tok.Lexema)
				if s == nil {
					return nil, errors.New("Stack error: Not declared ID ->" + expr.Tok.Lexema)
				}
			}

			if err != nil {
				return nil, err
			}
			nasig.Expr = expr
		}

		_, err, isPC := p.match(fxlex.TokPC) // ';'
		if err != nil || !isPC {
			return nil, errors.New("Error: Missing Dot$Comma in Asignation")
		}
		nsent.Asig = nasig
	case fxlex.TokLPar:
		nfcall := NewFcall() // Create FUNCTIONCALL AST
		nfcall.C = "FunctionCall"
		nfcall.depth = depth + 1
		nfcall.Id = r.Lexema

		t, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Missing Left Parenthesis in FunctionCall")
		}

		var args []gixsymb.Symb
		t, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Wrong argument in FunctionCall")
		}
		args = append(args, gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Argument"})

		args, err = p.Args(args)
		if err != nil {
			return nil, err
		}
		nfcall.Args = args

		_, err, isRPar := p.match(fxlex.TokRPar) // ')'
		if err != nil || !isRPar {
			return nil, errors.New("Error: Missing Right Parenthesis in FunctionCall")
		}
		_, err, isPC := p.match(fxlex.TokPC) // ';'
		if err != nil || !isPC {
			return nil, errors.New("Error: Missing Dot&Comma in FunctionCall")
		}

		nsent.Fcall = nfcall
	}
	return nsent, nil
}

// <INIT>	::=	<ATOM>	':='	<ATOM>	';'
func (p *Parser) Init(depth int) (asig *Init, err error) {
	p.pushTrace("Init")
	defer p.popTrace()

	ninit := NewInit() // Create INIT AST
	ninit.C = "Init"
	ninit.depth = depth + 1

	t, err, isID := p.match(fxlex.TokId) // Id
	if err != nil || !isID {
		return nil, errors.New("Error: Wrong ID in Init")
	}
	ninit.Id = gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Var"}

	s, err := p.stkEnv.NewSymb(t.Lexema, gixsymb.SVar)
	if err != nil || s == nil {
		return nil, errors.New("Stack error: ID already defined -> " + t.Lexema)
	}

	_, err, isAsig := p.match(fxlex.TokAsig) // INIT
	if err != nil || !isAsig {
		return nil, errors.New("Error: Missing AsignationToken in Init")
	}

	fexpr, err := p.Expr(1) // ATOM
	if err != nil {
		return nil, err
	}
	ninit.Expr = fexpr

	_, err, isPC := p.match(fxlex.TokPC) // ';'
	if err != nil || !isPC {
		return nil, errors.New("Error: Missing Dot&Comma in Init")
	}
	return ninit, nil
}

//	<ARGS>	::=	','	<ATOM>	<ARGS>	|
//							Empty
func (p *Parser) Args(args []gixsymb.Symb) (s []gixsymb.Symb, err error) {
	p.pushTrace("Args")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return s, err
	}
	switch t.Type {
	case fxlex.TokComma:
		_, err := p.l.Lex()
		if err != nil {
			return s, errors.New("Error: Missing Comma in Args")
		}

		t, err = p.l.Lex()
		if err != nil {
			return nil, errors.New("Error: Wrong argument in Args")
		}
		args = append(args, gixsymb.Symb{t.Lexema, 0, int(t.Type), float64(t.Valor), "Argument"})

		args, err := p.Args(args)
		if err != nil {
			return s, err
		}
		return args, nil
	default:
		return args, err
	}
}

// <ATOM>	::=	int	&	bool	&	Id	|	Empty
func (p *Parser) Atom() (atom string, err error) {
	p.pushTrace("Atom")
	defer p.popTrace()

	t, err := p.l.Peek()
	if err != nil {
		return "", err
	}
	switch t.Type {
	//<ATOM>	::=	int	&	Id
	case fxlex.TokInt, fxlex.TokId:
		t, err := p.l.Lex()
		if err != nil {
			return "", errors.New("Error: Wrong terminal in Atom")
		}

		if t.Type == fxlex.TokId {
			fmt.Println("TokType: ", t.Type, "id: ", t.Lexema)
			id := strings.Split(t.Lexema, ".")
			var s *gixsymb.Symb
			if len(id) == 1 {
				s = p.stkEnv.GetSymb(t.Lexema)
			} else {
				s = p.stkEnv.GetSymb(id[0])
			}
			if s == nil {
				return "", errors.New("Stack error: ID not declared -> " + t.Lexema)
			}
		}

		return t.Lexema, nil
		//<ATOM>	::=	'['	TokInt',	TokInt	']'
	case fxlex.TokLCorch:
		t, err := p.l.Lex()
		if err != nil {
			return "", errors.New("Error: Missing Left Bracket in Atom for Coord")
		}
		s := t.Lexema

		t, err, isTokInt := p.match(fxlex.TokInt) // TokInt
		if err != nil || !isTokInt {
			return "", errors.New("Error: Wrong Tokint in Atom for Coord")
		}
		s += t.Lexema

		t, err, isComma := p.match(fxlex.TokComma) // ','
		if err != nil || !isComma {
			return "", errors.New("Error: Missing Comma in Atom for Coord")
		}
		s += t.Lexema

		t, err, isTokInt = p.match(fxlex.TokInt) // TokInt
		if err != nil || !isTokInt {
			return "", errors.New("Error: Wrong Tokint in Atom for Coord")
		}
		s += t.Lexema

		t, err, isRCorch := p.match(fxlex.TokRCorch) // ']'
		if err != nil || !isRCorch {
			return "", errors.New("Error: Missing Right Bracket in Atom for Coord")
		}
		s += t.Lexema

		return s, err
	//<ATOM>	::=	Empty
	default:
		return "", errors.New("Syntax Error")
	}

}
