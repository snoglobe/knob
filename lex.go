package main

const (
	LeftParen uint8 = iota
	RightParen
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket
	Semicolon
	Quote
	Number
	Symbol
	Eof
)

type Lexer struct {
	Current int
	Source  string
}

type Token struct {
	Type uint8
	Val  string
}

func (l *Lexer) Advance() uint8 {
	if l.Current >= len(l.Source) {
		return 0
	}
	l.Current++
	return l.Source[l.Current-1]
}

func (l *Lexer) Peek() uint8 {
	if l.Current >= len(l.Source) {
		return 0
	}
	return l.Source[l.Current]
}

func (l *Lexer) ScanToken() Token {
start:
	t := Token{}
	c := l.Advance()

	switch c {
	case '(':
		t.Type = LeftParen
		t.Val = "("
	case ')':
		t.Type = RightParen
		t.Val = ")"
	case '{':
		t.Type = LeftBrace
		t.Val = "{"
	case '}':
		t.Type = RightBrace
		t.Val = "}"
	case '[':
		t.Type = LeftBracket
		t.Val = "["
	case ']':
		t.Type = RightBracket
		t.Val = "]"
	case ';':
		t.Type = Semicolon
		t.Val = ";"
	case '`':
		t.Type = Quote
		t.Val = "`"
	case ' ', '\t', '\r', '\n':
		goto start
	case 0:
		t.Type = Eof
		t.Val = ""
	default:
		if c >= '0' && c <= '9' {
			t.Type = Number
			t.Val = string(c)
			for l.Peek() >= '0' && l.Peek() <= '9' {
				t.Val += string(l.Advance())
			}
		} else {
			t.Type = Symbol
			t.Val = string(c)
			for l.Peek() != ' ' && l.Peek() != '\t' && l.Peek() != '\r' && l.Peek() != '\n' &&
				l.Peek() != '(' && l.Peek() != ')' && l.Peek() != '{' && l.Peek() != '}' &&
				l.Peek() != '[' && l.Peek() != ']' && l.Peek() != ';' && l.Peek() != '\'' && l.Peek() != 0 {
				t.Val += string(l.Advance())
			}
		}
	}
	return t
}
