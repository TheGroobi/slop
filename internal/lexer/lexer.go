package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TokenType = int

const (
	TOKEN_IDENT = iota
	TOKEN_DOT
	TOKEN_LBRACKET
	TOKEN_RBRACKET
	TOKEN_STRING
	TOKEN_NEWLINE
	TOKEN_EOF
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

type Lexer struct {
	input   *bufio.Reader
	current rune
	eof     bool
	line    int
}

func NewLexer(r io.Reader) *Lexer {
	l := &Lexer{
		input: bufio.NewReader(r),
		line:  1,
	}

	l.advance()
	return l
}

func (l *Lexer) NextToken() Token {
	for l.current == ' ' || l.current == '\t' || l.current == '\r' {
		l.advance()
	}

	if l.eof {
		return Token{Type: TOKEN_EOF, Literal: "", Line: l.line}
	}

	switch l.current {
	case '#':
		for l.current != '\n' && !l.eof {
			l.advance()
		}
	case '.':
		l.advance()
		return Token{Type: TOKEN_DOT, Literal: ".", Line: l.line}
	case '[':
		l.advance()
		return Token{Type: TOKEN_LBRACKET, Literal: "[", Line: l.line}
	case ']':
		l.advance()
		return Token{Type: TOKEN_RBRACKET, Literal: "]", Line: l.line}
	case '\n':
		l.advance()
		l.line++
		return Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: l.line}
	case '"':
		str := l.readString()
		return Token{Type: TOKEN_STRING, Literal: str, Line: l.line}
	default:
		if isLetter(l.current) {
			ident := l.readIdent()
			return Token{Type: TOKEN_IDENT, Literal: ident, Line: l.line}
		}
	}

	ch := l.current
	l.advance()
	fmt.Printf("warning: unexpected character '%c' on line %d\n", ch, l.line)
	return l.NextToken()
}

func (l *Lexer) advance() {
	ch, _, err := l.input.ReadRune()
	if err != nil {
		l.eof = true
		l.current = 0
	} else {
		l.current = ch
	}

}

func (l *Lexer) readIdent() string {
	var sb strings.Builder
	for isLetter(l.current) {
		sb.WriteRune(l.current)
		l.advance()
	}

	return sb.String()
}

func (l *Lexer) readString() string {
	var sb strings.Builder
	l.advance() // skip initial string quote

	for l.current != '"' && !l.eof {
		sb.WriteRune(l.current)
		l.advance()
	}

	l.advance() // skip last string quote

	return sb.String()
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
