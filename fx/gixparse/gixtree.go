package gixparse

import (
	"fmt"
	"fx/fxlex"
	"fx/gixsymb"
)

type Program struct {
	C        string
	Auxiliar []*Auxiliar
	depth    int
}

type Auxiliar struct {
	C           string
	Function    []*Function
	Declaration []*Declaration
	depth       int
}

type Declaration struct {
	C      string
	Id     gixsymb.Symb
	Params []string
	depth  int
}

type Function struct {
	C        string
	Id       string
	Params   []string
	Sentence []*Sentence
	depth    int
}

type Sentence struct {
	C     string
	Iter  *Iter
	Fcall *Fcall
	Sif   *Sif
	Selse *Selse
	Decl  gixsymb.Symb
	Asig  *Asig
	depth int
}

type Sif struct {
	C        string
	Id       string
	Expr     *Expr
	Sentence []*Sentence
	depth    int
}

type Selse struct {
	C        string
	Id       string
	Sentence []*Sentence
	depth    int
}

type Iter struct {
	C        string
	Id       string
	Init     *Init
	Atom1    gixsymb.Symb
	Atom2    gixsymb.Symb
	Sentence []*Sentence
	depth    int
}

type Fcall struct {
	C     string
	Id    string
	Args  []gixsymb.Symb
	depth int
}

type Init struct {
	C     string
	Id    gixsymb.Symb
	Expr  *Expr
	depth int
}

type Asig struct {
	C     string
	Id    gixsymb.Symb
	Expr  *Expr
	Coord string
	depth int
}

type Expr struct {
	C      string
	Tok    fxlex.Token
	ERight *Expr
	ELeft  *Expr
}

//type Atom struct {
//	tok    fxlex.Token
//	ERight *Atom
//	ELeft  *Atom
//}

func NewProgram() (prog *Program) {
	return &Program{Auxiliar: nil}
}

func NewAuxiliar() (prog *Auxiliar) {
	return &Auxiliar{Function: nil, Declaration: nil}
}

func NewDeclaration() (decl *Declaration) {
	return &Declaration{}
}

func NewFunction() (function *Function) {
	return &Function{Params: nil, Sentence: nil}
}

func NewSentence() (sentence *Sentence) {
	return &Sentence{Iter: nil, Fcall: nil, Sif: nil, Selse: nil}
}

func NewIter() (iter *Iter) {
	return &Iter{Init: nil, Sentence: nil}
}

func NewSif() (sif *Sif) {
	return &Sif{Expr: nil, Sentence: nil}
}

func NewSelse() (selse *Selse) {
	return &Selse{Sentence: nil}
}

func NewFcall() (fcall *Fcall) {
	return &Fcall{}
}

func NewInit() (Asig *Init) {
	return &Init{Expr: nil}
}

func NewAsig() (asig *Asig) {
	return &Asig{Expr: nil}
}

func NewExpr(tok fxlex.Token) (expr *Expr) {
	return &Expr{Tok: tok}
}

//func NewAtom(tok fxlex.Token) (Atom *Atom) {
//	return &Atom{tok: tok}
//}

// Para depurar
func (p *Program) PrintProg() string {
	if p == nil {
		return "nil"
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %v\n", p.C, p.depth, p.Auxiliar)

	for _, f := range p.Auxiliar {
		if f != nil {
			s += "	" + f.printAux()
		}
	}
	return s
}

func (a *Auxiliar) printAux() string {
	if a == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %v %v\n", a.C, a.depth, a.Function, a.Declaration)

	for _, t := range a.Declaration {
		if t != nil {
			s += "		" + t.printDecl()
		}
	}

	for _, t := range a.Function {
		if t != nil {
			s += "		" + t.printFunc()
		}
	}

	return s
}

func (d *Declaration) printDecl() string {
	if d == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %s\n", d.C, d.depth, d.Id.Name, d.Params)
	return s
}

func (f *Function) printFunc() string {
	if f == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v %v\n", f.C, f.depth, f.Id, f.Params, f.Sentence)

	for _, t := range f.Sentence {
		if t != nil {
			s += "		" + t.printSentence()
		}
	}
	return s
}

func (s *Sentence) printSentence() string {
	if s == nil {
		return ""
	}
	c := fmt.Sprintf("Node: %s, Depth: %d, Tree: %v %v %v %v %v %v\n", s.C, s.depth, s.Fcall, s.Asig, s.Iter, s.Sif, s.Selse, s.Decl)
	if s.Iter != nil {
		c += "			" + s.Iter.printIter()
	} else if s.Fcall != nil {
		c += "			" + s.Fcall.printFcall()
	} else if s.Asig != nil {
		c += "			" + s.Asig.printAsig()
	} else if s.Sif != nil {
		c += "			" + s.Sif.printSif()
	} else if s.Selse != nil {
		c += "			" + s.Selse.printSelse()
	}
	return c
}

func (i *Iter) printIter() string {
	if i == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v %s %s %v\n", i.C, i.depth, i.Id, i.Init, i.Atom1.Name, i.Atom2.Name, i.Sentence)

	if i.Init != nil {
		s += "				" + i.Init.printInit()
	}

	for _, t := range i.Sentence {
		if t != nil {
			s += "				" + t.printSentence()
		}
	}
	return s
}

func (i *Sif) printSif() string {
	if i == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v %v\n", i.C, i.depth, i.Id, i.Expr, i.Sentence)

	for _, t := range i.Sentence {
		if t != nil {
			s += "				" + t.printSentence()
		}
	}
	return s
}

func (i *Selse) printSelse() string {
	if i == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v\n", i.C, i.depth, i.Id, i.Sentence)

	for _, t := range i.Sentence {
		if t != nil {
			s += "				" + t.printSentence()
		}
	}
	return s
}

func (i *Init) printInit() string {
	if i == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v\n", i.C, i.depth, i.Id.Name, i.Expr)
	return s
}

func (a *Asig) printAsig() string {
	if a == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: Id_Name: %s, Coord: %s, Expr: %v\n", a.C, a.depth, a.Id.Name, a.Coord, a.Expr.printExpr())
	return s
}

func (f *Fcall) printFcall() string {
	if f == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Depth: %d, Tree: %s %v\n", f.C, f.depth, f.Id, f.Args)
	return s
}

func (e *Expr) printExpr() string {
	if e == nil {
		return ""
	}
	s := fmt.Sprintf("Node: %s, Token: %v, Tree: %v %v.", e.C, e.Tok, e.ELeft, e.ELeft)
	return s
}
