package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		run(os.Args[1])
	} else {
		repl()
	}
}

func run(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	text := string(file)
	lexer := Lexer{
		Current: 0,
		Source:  text,
	}
	parser := NewParser(&lexer)
	Eval(parser.Program())
}

func repl() {
	for {
		fmt.Printf(">>> ")
		input, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		lexer := Lexer{
			Current: 0,
			Source:  string(input),
		}
		parser := NewParser(&lexer)
		v := Eval(parser.Program())
		if v != nil {
			fmt.Println(v)
		}
	}
}
