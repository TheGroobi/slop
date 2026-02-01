package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type TokenType = int

const (
	TOKEN_IDENT        = iota // keyword
	TOKEN_DOT                 // .
	TOKEN_AT                  // @
	TOKEN_TASK_START          // {
	TOKEN_TASK_END            // }
	TOKEN_LBRACKET            // [
	TOKEN_RBRACKET            // ]
	TOKEN_STRING              // "
	TOKEN_NEWLINE             // \n
	TOKEN_EOF                 // EOF
	TOKEN_DOUBLE_COLON        // ::
	TOKEN_INTERP              // $var - variable reference / task reference
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
		line:  0,
	}

	fmt.Println("✔ Slop lexer initialized!")

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
		return l.NextToken()
	case '.':
		l.advance()
		return Token{Type: TOKEN_DOT, Literal: ".", Line: l.line}
	case ':':
		l.advance()
		if l.current == ':' {
			l.advance()
			return Token{Type: TOKEN_DOUBLE_COLON, Literal: "::", Line: l.line}
		}
	case '@':
		l.advance()
		return Token{Type: TOKEN_AT, Literal: "@", Line: l.line}
	case '{':
		l.advance()
		return Token{Type: TOKEN_TASK_START, Literal: "{", Line: l.line}
	case '}':
		l.advance()
		return Token{Type: TOKEN_TASK_END, Literal: "}", Line: l.line}
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
		str := l.readUntil('"')
		return Token{Type: TOKEN_STRING, Literal: str, Line: l.line}
	case '$':
		l.advance()
		ident := l.readDottedIdent()
		return Token{Type: TOKEN_INTERP, Literal: ident, Line: l.line}
	default:
		if isValidRune(l.current) {
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
	for isValidRune(l.current) {
		sb.WriteRune(l.current)
		l.advance()
	}

	return sb.String()
}

func (l *Lexer) readDottedIdent() string {
	var sb strings.Builder
	for isValidRune(l.current) || l.current == '.' {
		sb.WriteRune(l.current)
		l.advance()
	}

	return sb.String()
}

func (l *Lexer) readUntil(end rune) string {
	var sb strings.Builder
	l.advance() // skip initializer

	for l.current != end && !l.eof {
		sb.WriteRune(l.current)
		l.advance()
	}

	l.advance() // skip the finalizer

	return sb.String()
}

func isValidRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' || ch == '-'
}

func (l *Lexer) Lex() []Token {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == TOKEN_EOF {
			break
		}
	}

	return tokens
}
