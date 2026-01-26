package main

import (
	"bufio"
	"fmt"
	"io"
)

type Parser struct {
	Buf *bufio.Reader
}

func main() {
	var r io.Reader
	parser := newParser(r)

	if parser != nil {
		fmt.Println("Slop has been inited correctly")
	}
}

func newParser(r io.Reader) *Parser {
	return &Parser{
		Buf: bufio.NewReader(r),
	}
}
