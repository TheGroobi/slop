package parser

import (
	"fmt"

	"github.com/thegroobi/slop/internal/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		pos:     0,
		current: tokens[0],
	}
}

func (p *Parser) Parse() (map[string]string, error) {

	m := make(map[string]string)

	for p.current.Type != lexer.TOKEN_EOF {
		if p.current.Type == lexer.TOKEN_NEWLINE {
			p.advance()
			continue
		}

		k, v, err := p.parseDeclaration()
		if err != nil || k == "" || v == "" {
			return nil, err
		}

		m[k] = v

	}

	return m, nil
}

func (p *Parser) advance() {
	p.pos++
	p.current = p.tokens[p.pos]
}

func (p *Parser) parseLine() {}

func (p *Parser) parseDeclaration() (string, string, error) {
	// expect: IDENT DOT IDENT LBRACKET STRING RBRACKET

	if p.current.Type != lexer.TOKEN_IDENT {
		return "", "", fmt.Errorf("line %d: expected identifier", p.current.Line)
	}
	key := p.current.Literal
	p.advance()

	if p.current.Type != lexer.TOKEN_DOT {
		return "", "", fmt.Errorf("line %d: expected '.'", p.current.Line)
	}
	p.advance()

	if p.current.Type != lexer.TOKEN_IDENT {
		return "", "", fmt.Errorf("line %d: expected identifier after '.'", p.current.Line)
	}
	key = key + "." + p.current.Literal
	p.advance()

	if p.current.Type != lexer.TOKEN_LBRACKET {
		return "", "", fmt.Errorf("line %d: expected '['", p.current.Line)
	}
	p.advance()

	if p.current.Type != lexer.TOKEN_STRING {
		return "", "", fmt.Errorf("line %d: expected string", p.current.Line)
	}
	value := p.current.Literal
	p.advance()

	if p.current.Type != lexer.TOKEN_RBRACKET {
		return "", "", fmt.Errorf("line %d: expected ']'", p.current.Line)
	}
	p.advance()

	return key, value, nil
}
