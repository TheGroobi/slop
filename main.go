package main

import (
	"fmt"
	"os"

	"github.com/thegroobi/slop/internal/lexer"
	"github.com/thegroobi/slop/internal/parser"
)

const MIN_ARGS = 1
const SLOPFILE_NAME = "Slopfile"

func main() {
	// args := os.Args[1:]
	// err := validateArgs(args)
	// if err != nil {
	// 	fmt.Printf("%v: expected %d got %d\n", err, MIN_ARGS, len(args))
	// 	return
	// }
	//
	file, err := os.Open(SLOPFILE_NAME)
	if err != nil {
		fmt.Println("Could not read Slopfile", err)
		return
	}
	defer file.Close()

	fmt.Println("Slopfile loaded 🤤")

	l := lexer.NewLexer(file)
	if l != nil {
		fmt.Println("Slop lexer initialized!")
	}

	var tokens []lexer.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == lexer.TOKEN_EOF {
			break
		}
	}

	parser := parser.NewParser(tokens)
	if parser != nil {
		fmt.Println("Slop parser initialized!")
	}

	m, err := parser.Parse()

	if err != nil {
		fmt.Println("Parser failed to parse:", err)
	}

	fmt.Println(m)
}

// func validateArgs(args []string) error {
// 	if len(args) < MIN_ARGS {
// 		return errors.New("Invalid number of args")
// 	}
//
// 	return nil
// }
