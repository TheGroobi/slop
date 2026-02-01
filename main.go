package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegroobi/slop/internal/lexer"
	"github.com/thegroobi/slop/internal/parser"
	"github.com/thegroobi/slop/internal/slopfile"
)

const SLOPFILE_NAME = "Slopfile"

func main() {
	file, err := os.Open(SLOPFILE_NAME)
	if err != nil {
		fmt.Println("✘ Could not read Slopfile", err)
		return
	}
	defer file.Close()

	Init()

	l := lexer.NewLexer(file)
	tokens := l.Lex()
	if len(tokens) == 0 {
		fmt.Println("✘ Config empty: returning early")
		return
	}

	parser := parser.NewParser(tokens)

	slop := slopfile.NewSlopfile()

	slop, err = parser.Parse(slop)
	if err != nil {
		fmt.Println("✘ Parser failed:", err)
		return
	}

	fmt.Println("----------------------------------------------")

	if err = slop.Run(os.Args[1:]); err != nil {
		fmt.Println("✘ Run failed:", err)
		return
	}

	if slop.DirectivesRun == 0 {
		fmt.Println("𝒊 Nothing was run - provide some directives in the Slopfile")
	}
}

func Init() {
	godotenv.Load()

	fmt.Println("✔ Slopfile loaded 🤤")
}
