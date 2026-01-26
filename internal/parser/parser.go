package parser

import (
	"bufio"
	"io"
)

type Parser struct {
	Buf *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		Buf: bufio.NewReader(r),
	}
}
