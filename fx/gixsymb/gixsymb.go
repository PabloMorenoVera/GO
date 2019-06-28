package gixsymb

import (
	"errors"
)

const (
	SNone  = iota
	SConst //Constant
	SType  //type definition
	SVar   //object definition
	SFunc  //function
)

type Symb struct {
	Name  string
	SType int
	//DataType *Type
	TokKind  int
	FloatVal float64
	FName    string
}

type Env map[string]*Symb

type StkEnv []Env

func (envs *StkEnv) PushEnv(s string) {
	env := Env{}
	*envs = append(*envs, env)
	//fmt.Println("Env Created["+s+"]: ", envs)
}

func (envs *StkEnv) PopEnv(s string) {
	eS := *envs

	//fmt.Println("Env Deleted["+s+"]: ", envs)
	if len(eS) == 1 {
		panic("Cannot pop builtin")
	}
	*envs = eS[:len(eS)-1]
}

func (envs *StkEnv) NewSymb(name string, sType int) (s *Symb, err error) {
	eS := *envs
	s = &Symb{SType: sType}
	e := eS[len(eS)-1]

	if _, ok := e[name]; ok {
		return nil, errors.New("Already declared symb -> " + name)
	}
	e[name] = s
	return s, nil
}

func (envs *StkEnv) GetSymb(name string) (s *Symb) {
	eS := *envs

	for i := len(eS) - 1; i >= 0; i-- {
		if s, ok := eS[i][name]; ok {
			return s
		}
	}
	return nil
}
