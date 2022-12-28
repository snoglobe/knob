package main

import "fmt"

var variables map[string]any

func init() {
	variables = map[string]any{
		"print": func(args []Atom) any {
			for _, a := range args {
				fmt.Printf("%v", Eval(a))
			}
			fmt.Printf("\n")
			return nil
		},
		"do": func(args []Atom) any {
			var v any
			for _, a := range args {
				v = Eval(a.Val.(Atom))
			}
			return v
		},
		"set": func(args []Atom) any {
			if len(args) != 2 {
				panic("set: expected 2 arguments")
			}
			name := args[0]
			if name.Type != Sym {
				panic("set: expected symbol as first argument")
			}
			variables[name.Val.(string)] = Eval(args[1])
			return nil
		},
		"+": func(args []Atom) any {
			var v int
			for _, a := range args {
				v += Eval(a).(int)
			}
			return v
		},
		"-": func(args []Atom) any {
			var v int
			for i, a := range args {
				if i == 0 {
					v = Eval(a).(int)
				} else {
					v -= Eval(a).(int)
				}
			}
			return v
		},
		"*": func(args []Atom) any {
			var v = 1
			for _, a := range args {
				v *= Eval(a).(int)
			}
			return v
		},
		"/": func(args []Atom) any {
			var v int
			for i, a := range args {
				if i == 0 {
					v = Eval(a).(int)
				} else {
					v /= Eval(a).(int)
				}
			}
			return v
		},
		"=": func(args []Atom) any {
			if len(args) != 2 {
				panic("=: expected 2 arguments")
			}
			return Eval(args[0]) == Eval(args[1])
		},
		"!=": func(args []Atom) any {
			if len(args) != 2 {
				panic("!=: expected 2 arguments")
			}
			return Eval(args[0]) != Eval(args[1])
		},
		"<": func(args []Atom) any {
			if len(args) != 2 {
				panic("<: expected 2 arguments")
			}
			return Eval(args[0]).(int) < Eval(args[1]).(int)
		},
		"<=": func(args []Atom) any {
			if len(args) != 2 {
				panic("<=: expected 2 arguments")
			}
			return Eval(args[0]).(int) <= Eval(args[1]).(int)
		},
		">": func(args []Atom) any {
			if len(args) != 2 {
				panic(">: expected 2 arguments")
			}
			return Eval(args[0]).(int) > Eval(args[1]).(int)
		},
		">=": func(args []Atom) any {
			if len(args) != 2 {
				panic(">=: expected 2 arguments")
			}
			return Eval(args[0]).(int) >= Eval(args[1]).(int)
		},
		"if": func(args []Atom) any {
			if len(args) != 3 {
				panic("if: expected 3 arguments")
			}
			if Eval(args[0]).(bool) {
				return Eval(args[1])
			}
			return Eval(args[2])
		},
		":true":  true,
		":false": false,
		"and": func(args []Atom) any {
			for _, a := range args {
				if !Eval(a).(bool) {
					return false
				}
			}
			return true
		},
		"or": func(args []Atom) any {
			for _, a := range args {
				if Eval(a).(bool) {
					return true
				}
			}
			return false
		},
		"not": func(args []Atom) any {
			if len(args) != 1 {
				panic("not: expected 1 argument")
			}
			return !Eval(args[0]).(bool)
		},
		"lambda": func(args []Atom) any {
			if len(args) != 2 {
				panic("lambda: expected 2 arguments")
			}
			var params []Atom
			if len(args[0].Val.([]Atom)) > 0 {
				zeroArg := args[0]
				if zeroArg.Type != List {
					panic("lambda: expected list as first argument")
				}
				params = zeroArg.Val.([]Atom)
			}
			body := args[1]
			return func(args []Atom) any {
				if params != nil && len(params) != len(args) {
					panic("lambda: expected " + fmt.Sprint(len(params)) + " arguments")
				}
				tempVars := variables
				variables = map[string]any{}
				for k, v := range tempVars {
					variables[k] = v
				}
				for i, p := range params {
					variables[p.Val.(string)] = Eval(args[i])
				}
				v := Eval(body)
				variables = tempVars
				return v
			}
		},
		"macro": func(args []Atom) any {
			if len(args) != 2 {
				panic("lambda: expected 2 arguments")
			}
			var params []Atom
			if len(args[0].Val.([]Atom)) > 0 {
				zeroArg := args[0]
				if zeroArg.Type != List {
					panic("lambda: expected list as first argument")
				}
				params = zeroArg.Val.([]Atom)
			}
			body := args[1]
			return func(args []Atom) any {
				if params != nil && len(params) != len(args) {
					panic("lambda: expected " + fmt.Sprint(len(params)) + " arguments")
				}
				tempVars := variables
				variables = map[string]any{}
				for k, v := range tempVars {
					variables[k] = v
				}
				for i, p := range params {
					variables[p.Val.(string)] = args[i]
				}
				v := Eval(Eval(body).(Atom))
				variables = tempVars
				return v
			}
		},
		"eval": func(args []Atom) any {
			if len(args) != 1 {
				panic("eval: expected 1 argument")
			}
			return Eval(args[0])
		},
		"incr": func(args []Atom) any {
			if len(args) != 1 {
				panic("incr: expected 1 argument")
			}
			variables[args[0].Val.(string)] = variables[args[0].Val.(string)].(int) + 1
			return variables[args[0].Val.(string)]
		},
		"decr": func(args []Atom) any {
			if len(args) != 1 {
				panic("decr: expected 1 argument")
			}
			variables[args[0].Val.(string)] = variables[args[0].Val.(string)].(int) - 1
			return variables[args[0].Val.(string)]
		},
		"for": func(args []Atom) any {
			if len(args) != 4 {
				panic("for: expected 4 arguments")
			}
			init := args[0]
			cond := args[1]
			incr := args[2]
			body := args[3]
			Eval(init)
			for Eval(cond).(bool) {
				Eval(body)
				Eval(incr)
			}
			return nil
		},
	}
}

func Eval(a Atom) any {
	switch a.Type {
	case Num:
		return a.Val.(int)
	case Quot:
		return a.Val.(Atom)
	case List:
		l := a.Val.([]Atom)
		if len(l) == 0 {
			return nil
		}
		fnc := Eval(l[0])
		args := make([]Atom, 0)
		for _, a := range l[1:] {
			args = append(args, a)
		}
		return (fnc.(func([]Atom) any))(args)
	case Sym:
		if v, ok := variables[a.Val.(string)]; ok {
			return v
		}
		panic("eval: undefined variable " + a.Val.(string))
	}
	panic("eval: unhandled atom type " + a.Type.String())
}
