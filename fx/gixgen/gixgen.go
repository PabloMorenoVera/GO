package gixgen

import (
	"fmt"
	"fx/fxlex"
	"fx/gixparse"
	"fx/gixsymb"
	"strconv"
	"strings"
)

type Stack struct {
	Name string
	Args []gixsymb.Symb
	Vars []gixsymb.Symb
	Cont int
}

type LabelStack struct {
	items []Stack
}

func NewStack() *LabelStack {
	lstk := &LabelStack{}
	return lstk
}

func (lstk *LabelStack) pushLabel(item Stack) {
	lstk.items = append(lstk.items, item)
}

func (lstk *LabelStack) popLabel() (item Stack) {
	l := len(lstk.items)
	item = lstk.items[l-1]
	lstk.items = lstk.items[:l-1]
	return item
}

func (lstk *LabelStack) ProgMain(node *gixparse.Program) *LabelStack {
	ifbool := true
	item := Stack{"main", nil, nil, 0}
	lstk.pushLabel(item)

	for _, a := range node.Auxiliar {
		// Search for global declarations
		for _, d := range a.Declaration {
			lstk = lstk.GlobalDeclGen(d, lstk)
		}
		// Search for "main" function
		for _, f := range a.Function {
			if f.Id == "main" {
				fmt.Println("# main")
				for _, s := range f.Sentence {
					if s.Decl.Name != "" {
						lstk.DeclGen(s, lstk)
					} else if s.Iter != nil {
						lstk.IterGen(node, lstk, s.Iter)
					} else if s.Fcall != nil {
						fmt.Println("#", s.Fcall.Id)
						id, lstk := lstk.FcallGen(s, lstk)
						if id != "" {
							lstk.FuncGen(node, lstk, id)
						}
						fmt.Println("# main")
					} else if s.Asig != nil {
						lstk.AsigGen(s, lstk)
					} else if s.Sif != nil {
						ifbool, lstk = lstk.IfGen(node, s.Sif, lstk)
					} else if s.Selse != nil && !ifbool {
						lstk.ElseGen(node, s.Selse, lstk)
					}
				}
			}
		}
	}
	return lstk
}

// GLOBAL_DECLARATION code generation
func (lstk *LabelStack) GlobalDeclGen(s *gixparse.Declaration, stk *LabelStack) *LabelStack {
	item := stk.popLabel()
	item.Vars = append(item.Vars, s.Id)
	stk.pushLabel(item)

	return stk
}

// DECLARATION code generation
func (lstk *LabelStack) DeclGen(s *gixparse.Sentence, stk *LabelStack) *LabelStack { // Insert Var in stack
	item := stk.popLabel()
	item.Vars = append(item.Vars, s.Decl)
	stk.pushLabel(item)

	return stk
}

// ITER code generation
func (lstk *LabelStack) IterGen(node *gixparse.Program, stk *LabelStack, iter *gixparse.Iter) *LabelStack {
	fmt.Println("# ", iter.Id, "(", iter.Init.Id.Name, ":=", iter.Init.Expr.Tok.Lexema, ",", iter.Atom1.Name, ",", iter.Atom2.Name, ");")

	item := stk.popLabel()
	svar := gixsymb.Symb{iter.Init.Id.Name, 0, iter.Init.Id.TokKind, float64(iter.Init.Expr.Tok.Valor), iter.Init.Id.FName}
	item.Vars = append(item.Vars, svar)

	stk.pushLabel(item)

	var end int
	var step int

	for _, v := range item.Args { // Busco el valor en los argumentos
		if v.Name == iter.Atom1.Name {
			end = int(v.FloatVal)
		} else if v.Name == iter.Atom2.Name {
			step = int(v.FloatVal)
		}
	}

	for _, n := range item.Vars { // Busco el valor en las variables
		if n.Name == iter.Atom1.Name {
			end = int(n.FloatVal)
		} else if n.Name == iter.Atom2.Name {
			step = int(n.FloatVal)
		}
	}

	// Compruebo que sean literales
	if end == 0 {
		end, _ = strconv.Atoi(iter.Atom1.Name)
	}
	if step == 0 {
		step, _ = strconv.Atoi(iter.Atom2.Name)
	}

	for i := 0; i <= end; i += step {
		//Ejecuto el código de Iter
		for _, s := range iter.Sentence {
			if s.Decl.Name != "" {
				lstk.DeclGen(s, lstk)
			} else if s.Iter != nil {
				lstk.IterGen(node, lstk, s.Iter)
			} else if s.Fcall != nil {
				id, lstk := lstk.FcallGen(s, lstk)
				if id != "" {
					lstk.FuncGen(node, lstk, id)
				}
			} else if s.Asig != nil {
				lstk.AsigGen(s, lstk)
			} else if s.Sif != nil {
				lstk.IfGen(node, s.Sif, lstk)
			}
		}

		item := stk.popLabel()
		for i, n := range item.Vars {
			if n.Name == iter.Init.Id.Name {
				item.Vars[i].FloatVal = n.FloatVal + float64(step)
			}
		}
		stk.pushLabel(item)
	}

	return stk
}

// FUNCTIONCALL code generation
func (lstk *LabelStack) FcallGen(s *gixparse.Sentence, stk *LabelStack) (id string, st *LabelStack) {
	args := Stack{"Fcall", nil, nil, 0}
	item := stk.popLabel()

	if s.Fcall.Id == "rect" || s.Fcall.Id == "circle" {
		funcstring := fmt.Sprintf("%s ", s.Fcall.Id)
		for i, arg := range s.Fcall.Args { // Check all function arguments
			for _, val := range item.Vars { // Search the args in the stack
				if val.Name == (arg.Name+".x") || val.Name == (arg.Name+".y") || val.Name == (arg.Name+".z") { // Busco el campo entero en mis variables
					funcstring += fmt.Sprintf("%f ", val.FloatVal)
				}
			}
			if i != 0 { // Hay que quitar pp (se cambia por pp.x y pp.y)
				funcstring += fmt.Sprintf("%f ", arg.FloatVal)
			}
		}

		fstring := "  #" + s.Fcall.Id + "(" + s.Fcall.Args[0].Name
		for i, _ := range s.Fcall.Args {
			if i >= 1 {
				fstring += "," + s.Fcall.Args[i].Name
			}
		}
		fstring += ");"

		fmt.Println(funcstring, fstring)
		stk.pushLabel(item)
		return "", stk
	} else {
		for _, val := range item.Vars { // Search the args in the stack
			for _, arg := range s.Fcall.Args {
				if val.Name == arg.Name+".x" || val.Name == arg.Name+".y" || val.Name == arg.Name+".z" { // Busco el campo entero en mis variables
					args.Args = append(args.Args, val)
				}
			}
		}

		for _, a := range s.Fcall.Args {
			args.Args = append(args.Args, a)
		}
		stk.pushLabel(item)
		stk.pushLabel(args)
		return s.Fcall.Id, stk
	}
}

// FUNCTION code generation
func (lstk *LabelStack) FuncGen(node *gixparse.Program, stk *LabelStack, id string) *LabelStack {
	for _, a := range node.Auxiliar {
		for _, f := range a.Function {
			if f.Id == id {
				for _, s := range f.Sentence {
					if s.Decl.Name != "" {
						lstk.DeclGen(s, lstk)
					} else if s.Iter != nil {
						lstk.IterGen(node, lstk, s.Iter)
					} else if s.Fcall != nil {
						id, lstk := lstk.FcallGen(s, lstk)
						lstk.FuncGen(node, lstk, id)
					} else if s.Asig != nil {
						lstk.AsigGen(s, lstk)
					} else if s.Sif != nil {
						lstk.IfGen(node, s.Sif, lstk)
					}
				}
			}
		}
	}
	stk.popLabel()
	return stk
}

// ASIGNATION code generation
func (lstk *LabelStack) AsigGen(s *gixparse.Sentence, stk *LabelStack) *LabelStack {
	updated := false

	item := stk.popLabel()
	for i, t := range item.Vars { // Check  for the Var
		if t.Name == s.Asig.Id.Name {
			if s.Asig.Expr == nil { // Coord
				v, _ := strconv.ParseFloat(strings.Split(s.Asig.Coord, ",")[0][1:], 64)
				item.Vars = append(item.Vars, gixsymb.Symb{s.Asig.Id.Name + ".x", 0, int(fxlex.TokInt), v, "Var"})

				v, _ = strconv.ParseFloat(strings.Split(s.Asig.Coord, ",")[1][:len(strings.Split(s.Asig.Coord, ",")[1])-1], 64)
				item.Vars = append(item.Vars, gixsymb.Symb{s.Asig.Id.Name + ".y", 0, int(fxlex.TokInt), v, "Var"})
			} else { // It's an EXPR with declared var -> update
				updated = true
				if s.Asig.Expr.Tok.Valor == 0 {
					val := lstk.ExprGen(s.Asig.Expr, item)
					item.Vars[i].FloatVal = val
				} else {
					item.Vars[i].FloatVal = float64(s.Asig.Expr.Tok.Valor)
				}
			}
		}
	}
	if !updated {
		for _, t := range item.Vars { // Check for the Var declaration
			id := strings.Split(s.Asig.Id.Name, ".")
			if len(id) == 2 && id[0] == t.Name {
				val := lstk.ExprGen(s.Asig.Expr, item)
				item.Vars = append(item.Vars, gixsymb.Symb{s.Asig.Id.Name, 0, int(fxlex.TokId), val, "Var"})
			}
		}
	}
	stk.pushLabel(item)

	return stk
}

// IF code generation
func (lstk *LabelStack) IfGen(node *gixparse.Program, s *gixparse.Sif, stk *LabelStack) (elsecheck bool, st *LabelStack) { // Comprobar que se cumple la expresion
	item := lstk.popLabel()
	_, result := lstk.ExprChk(s.Expr, item)
	lstk.pushLabel(item)
	if result { // Miramos si se cumple la igualdad
		for _, a := range s.Sentence {
			if a.Decl.Name != "" {
				lstk.DeclGen(a, lstk)
			} else if a.Iter != nil {
				lstk.IterGen(node, lstk, a.Iter)
			} else if a.Fcall != nil {
				id, lstk := lstk.FcallGen(a, lstk)
				if id != "" {
					lstk.FuncGen(node, lstk, id)
				}
			} else if a.Asig != nil {
				lstk.AsigGen(a, lstk)
			} else if a.Sif != nil {
				lstk.IfGen(node, a.Sif, lstk)
			}
		}
	} else { // EXPR Check is false so go to FALSE
		return false, stk
	}
	return true, stk
}

// ELSE code generation. This code is part of the IF code because it call its
func (lstk *LabelStack) ElseGen(node *gixparse.Program, s *gixparse.Selse, stk *LabelStack) *LabelStack { // Si no se cumple la expresion se llama aquí
	for _, a := range s.Sentence {
		if a.Decl.Name != "" {
			lstk.DeclGen(a, lstk)
		} else if a.Iter != nil {
			lstk.IterGen(node, lstk, a.Iter)
		} else if a.Fcall != nil {
			id, lstk := lstk.FcallGen(a, lstk)
			if id != "" {
				lstk.FuncGen(node, lstk, id)
			}
		} else if a.Asig != nil {
			lstk.AsigGen(a, lstk)
		} else if a.Sif != nil {
			lstk.IfGen(node, a.Sif, lstk)
		}
	}
	return stk
}

// EXPRESION code generation -> return a value
func (lstk *LabelStack) ExprGen(e *gixparse.Expr, item Stack) float64 {
	rV := 0.0
	lV := 0.0

	if e == nil {
		return 0
	}
	if e.ERight != nil {
		rV = lstk.ExprGen(e.ERight, item)
	}
	if e.ELeft != nil {
		lV = lstk.ExprGen(e.ELeft, item)
	}
	tok := e.Tok
	switch tok.Type {
	case fxlex.TokSum:
		return lV + rV
	case fxlex.TokMul:
		return lV * rV
	case fxlex.TokRest:
		return lV - rV
	case fxlex.TokBar:
		return lV / rV
	case fxlex.TokInt:
		return float64(tok.Valor)
	case fxlex.TokId:
		val := 0.0
		for _, t := range item.Args { // Search the var in the args stack
			if t.Name == tok.Lexema {
				val = t.FloatVal
			}
		}
		for _, t := range item.Vars { // Search var in the var stack
			if t.Name == tok.Lexema {
				val = t.FloatVal
			}
		}
		return val
	default:
		panic("Bad subtree")
	}
}

// EXPRESSION-CHECK code generation
func (lstk *LabelStack) ExprChk(e *gixparse.Expr, item Stack) (f float64, b bool) {
	rV := 0.0
	lV := 0.0
	var boolif bool = b

	if e.ELeft != nil {
		lV, boolif = lstk.ExprChk(e.ELeft, item)
		if e.Tok.Type == fxlex.TokOr && boolif {
			return 0.0, boolif
		}
	}
	if e == nil {
		return 0.0, false
	}
	if e.ERight != nil {
		rV, boolif = lstk.ExprChk(e.ERight, item)
		if e.Tok.Type == fxlex.TokOr && boolif {
			return 0.0, boolif
		}
	}

	tok := e.Tok
	switch tok.Type {
	case fxlex.TokMin:
		return 0.0, lV < rV
	case fxlex.TokMax:
		return 0.0, lV > rV
	case fxlex.TokOr:
		if boolif == true {
			return 0.0, boolif
		}
		return lV, boolif
	case fxlex.TokAnd:
		if boolif == false {
			return 0.0, false
		}
		return lV, boolif
	case fxlex.TokNot:
		return rV, !boolif
	case fxlex.TokOpInt:
		if e.Tok.Lexema == "<=" {
			return 0.0, lV <= rV
		} else if e.Tok.Lexema == ">=" {
			return 0.0, lV >= rV
		}
		panic("Bad integer operator")
	case fxlex.TokBool:
		if tok.Lexema == "True" {
			//fmt.Println("Expr:", e, "Boolif", boolif)
			boolif = true
			return 0.0, true
		} else {
			return 0.0, false
		}
	case fxlex.TokInt:
		return float64(tok.Valor), boolif
	case fxlex.TokId:
		val := 0.0
		for _, t := range item.Args { // Search the var in the args stack
			if t.Name == tok.Lexema {
				val = t.FloatVal
			}
		}
		for _, t := range item.Vars { // Search var in the var stack
			if t.Name == tok.Lexema {
				val = t.FloatVal
			}
		}
		return val, boolif
	default:
		panic("Bad subtree")
	}
}
