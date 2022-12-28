package main

import (
	"strconv"
)

type AtomType uint8

//go:generate stringer -type=AtomType
const (
	Num  AtomType = iota
	Quot AtomType = iota
	List AtomType = iota
	Sym  AtomType = iota
)

type Atom struct {
	Type AtomType
	Val  any
}

func (a Atom) String() string {
	switch a.Type {
	case Num:
		return strconv.Itoa(a.Val.(int))
	case Sym:
		return a.Val.(string)
	case Quot:
		return "'" + a.Val.(Atom).String()
	case List:
		items := a.Val.([]Atom)
		s := "("
		for i, item := range items {
			s += item.String()
			if i < len(items)-1 {
				s += " "
			}
		}
		s += ")"
		return s
	}
	panic("atom: unexpected type")
}

func NewSymbol(s string) Atom {
	return Atom{
		Type: Sym,
		Val:  s,
	}
}

func NewList(a ...Atom) Atom {
	return Atom{
		Type: List,
		Val:  a,
	}
}

type Parser struct {
	Current Token
	Lexer   *Lexer
}

func NewParser(l *Lexer) *Parser {
	return &Parser{
		Current: l.ScanToken(),
		Lexer:   l,
	}
}

func (p *Parser) Peek() uint8 {
	return p.Current.Type
}

func (p *Parser) Eat(Type uint8) Token {
	if p.Peek() == Type {
		t := p.Current
		p.Current = p.Lexer.ScanToken()
		return t
	}
	panic("parser: unexpected token")
}

func (p *Parser) Program() Atom {
	items := []Atom{}
	for p.Peek() != Eof {
		subItems := []Atom{}
		for p.Peek() != Eof && p.Peek() != Semicolon {
			subItems = append(subItems, p.Parse())
		}
		if p.Peek() == Semicolon {
			p.Eat(Semicolon)
		}
		if len(subItems) == 1 {
			items = append(items, Atom{Quot, subItems[0]})
		} else {
			items = append(items, Atom{Quot, NewList(subItems...)})
		}
	}
	items = append([]Atom{NewSymbol("do")}, items...)
	return NewList(items...)
}

func (p *Parser) Parse() Atom {
	switch p.Peek() {
	case LeftParen:
		p.Eat(LeftParen)
		items := []Atom{}
		for p.Peek() != RightParen {
			items = append(items, p.Parse())
		}
		p.Eat(RightParen)
		return Atom{
			Type: List,
			Val:  items,
		}
	case Quote:
		p.Eat(Quote)
		return Atom{
			Type: Quot,
			Val:  p.Parse(),
		}
	case LeftBracket:
		p.Eat(LeftBracket)
		a := p.Parse()
		b := p.Parse()
		c := p.Parse()
		p.Eat(RightBracket)
		return Atom{
			Type: List,
			Val:  []Atom{b, a, c},
		}
	case LeftBrace:
		p.Eat(LeftBrace)
		items := []Atom{}
		for p.Peek() != RightBrace {
			subItems := []Atom{}
			for p.Peek() != RightBrace && p.Peek() != Semicolon {
				subItems = append(subItems, p.Parse())
			}
			if p.Peek() == Semicolon {
				p.Eat(Semicolon)
			}
			if len(subItems) == 1 {
				items = append(items, Atom{Quot, subItems[0]})
			} else {
				items = append(items, Atom{Quot, NewList(subItems...)})
			}
		}
		p.Eat(RightBrace)
		items = append([]Atom{NewSymbol("do")}, items...)
		return NewList(items...)
	case Symbol:
		t := p.Eat(Symbol)
		return Atom{
			Type: Sym,
			Val:  t.Val,
		}
	case Number:
		t := p.Eat(Number)
		v, _ := strconv.Atoi(t.Val)
		return Atom{
			Type: Num,
			Val:  v,
		}
	}
	panic("parser: unexpected token")
}
