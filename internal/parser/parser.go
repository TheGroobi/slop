package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/thegroobi/slop/internal/actions"
	"github.com/thegroobi/slop/internal/lexer"
)

const MAX_INDENTATION = 3

type Parser struct {
	tokens  []lexer.Token
	pos     int
	current lexer.Token
}

func NewParser(tokens []lexer.Token) *Parser {
	p := &Parser{
		tokens: tokens,
		pos:    0,
	}

	if len(tokens) > 0 {
		p.current = tokens[0]
	}

	return p
}

func (p *Parser) Parse() (*Slopfile, error) {
	slop := NewSlopfile()
	task := ""

	for p.current.Type != lexer.TOKEN_EOF {
		if p.current.Type == lexer.TOKEN_NEWLINE {
			p.advance()
			continue
		} else if p.current.Type == lexer.TOKEN_TASK_END {
			p.advance()
			task = ""
		}

		dir, k, v, err := p.parseDeclaration()
		if err != nil || k == "" || (v == "" && dir != DIR_TASK) {
			return nil, err
		}

		if strings.HasPrefix(v, "$env") {
			envKey := v[5:]
			envVal := os.Getenv(envKey)
			if envVal == "" {
				return nil, fmt.Errorf("line %d: environmental variable not set: %s", p.current.Line, envKey)
			}
			v = envVal
		} else if strings.HasPrefix(v, "$") {
			varName := v[1:]
			res, ok := slop.Vars[varName]
			if ok {
				v = res
			} else {
				_, ok := slop.Tasks[varName]
				if !ok {
					return nil, fmt.Errorf("line %d: undefined: $%s", p.current.Line, varName)
				}
				v = varName
			}
		}

		switch dir {
		case DIR_CONFIG:
			slop.Config[k] = v
		case DIR_VAR:
			slop.Vars[k] = v
		case DIR_SOURCE:
			action, err := actions.ParseAction(k)
			if err != nil {
				return nil, err
			} else if action != actions.ACT_ENV {
				return nil, fmt.Errorf("line %d: unknown action for source directive - allowed actions are: env", p.current.Line)
			}
			if err = godotenv.Load(v); err != nil {
				return nil, err
			}

			fmt.Printf("✔ .env file loaded from %s\n", v)
		case DIR_TASK:
			task = k
		case DIR_RUN:
			action, err := actions.ParseAction(k)
			if err != nil {
				return nil, err
			}

			run := actions.Action{
				Action: action,
				Args:   v,
				Line:   p.current.Line - 1,
			}

			if task != "" {
				slop.Tasks[task] = append(slop.Tasks[task], run)
			} else {
				slop.Runs = append(slop.Runs, run)
			}
		}

	}

	fmt.Println("✔ Slopfile has been parsed successfully")
	return slop, nil
}

func (p *Parser) advance() {
	p.pos++
	if p.pos < len(p.tokens) {
		p.current = p.tokens[p.pos]
	}
}

func (p *Parser) parseDeclaration() (DirectiveType, string, string, error) {
	// expect: IDENT(DIRECTIVE) DOUBLE_COLON IDENT(ACTION) DOT IDENT LBRACKET STRING RBRACKET

	dir, error := ParseDirective(p.current.Literal)
	if error != nil {
		return -1, "", "", fmt.Errorf("line %d: unexpected directive, available are: run, var, config, @task", p.current.Line)
	}
	p.advance()

	// return early when start task definition
	if dir == DIR_TASK {
		key := p.current.Literal
		p.advance()

		if p.current.Type != lexer.TOKEN_TASK_START {
			return -1, "", "", fmt.Errorf("line %d: unexpected task declaration missing opening bracket: \"{\"", p.current.Line)
		}

		p.advance() // skip "{"
		return dir, key, "", nil
	}

	// ::
	if p.current.Type != lexer.TOKEN_DOUBLE_COLON {
		return -1, "", "", fmt.Errorf("line %d: expected '::'", p.current.Line)
	}
	p.advance()

	// key identifier (for run directives this is validated as an action later)
	if p.current.Type != lexer.TOKEN_IDENT {
		return -1, "", "", fmt.Errorf("line %d: expected identifier after '::'", p.current.Line)
	}
	key := p.current.Literal
	p.advance()

	depth := 1
	// indentation max 3 depth
	for p.current.Type == lexer.TOKEN_DOT && depth < MAX_INDENTATION {
		p.advance() // skip dot
		if p.current.Type != lexer.TOKEN_IDENT {
			return -1, "", "", fmt.Errorf("line %d: expected identifier after '.'", p.current.Line)
		}
		key = key + "." + p.current.Literal
		p.advance()
		depth++
	}

	// [ value ]
	if p.current.Type != lexer.TOKEN_LBRACKET {
		return -1, "", "", fmt.Errorf("line %d: expected '['", p.current.Line)
	}
	p.advance()

	if p.current.Type != lexer.TOKEN_STRING && p.current.Type != lexer.TOKEN_INTERP {
		return -1, "", "", fmt.Errorf("line %d: expected string or variable reference", p.current.Line)
	}
	value := p.current.Literal
	if p.current.Type == lexer.TOKEN_INTERP {
		value = "$" + value
	}
	p.advance()

	if p.current.Type != lexer.TOKEN_RBRACKET {
		return -1, "", "", fmt.Errorf("line %d: expected ']'", p.current.Line)
	}
	p.advance()

	return dir, key, value, nil
}
