package parser

import (
	"fmt"
)

type DirectiveType int

const (
	DIR_RUN DirectiveType = iota
	DIR_CONFIG
	DIR_VAR
	DIR_SOURCE
	DIR_TASK // special task definition
)

type Directive struct {
	Directive DirectiveType
	Args      string
}

func ParseDirective(s string) (DirectiveType, error) {
	switch s {
	case "run":
		return DIR_RUN, nil
	case "@":
		return DIR_TASK, nil
	case "source":
		return DIR_SOURCE, nil
	case "config":
		return DIR_CONFIG, nil
	case "var":
		return DIR_VAR, nil
	default:
		return -1, fmt.Errorf("unknown directive: %s", s)
	}
}
