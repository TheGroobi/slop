package parser

import (
	"fmt"

	"github.com/thegroobi/slop/internal/actions"
)

type DirectiveType int

const (
	DIR_RUN DirectiveType = iota
	DIR_CONFIG
	DIR_VAR
	DIR_SOURCE
)

type Slopfile struct {
	Config map[string]string // config::db.name["x"] → Config["db.name"] = "x"
	Vars   map[string]string // var::seed.rbac["x"] → Vars["seed.rbac"] = "x"
	Runs   []actions.Action  // run::seed["x"] → append to Runs
}

type Directive struct {
	Directive DirectiveType
	Args      string
}

func NewSlopfile() *Slopfile {
	return &Slopfile{
		Config: make(map[string]string),
		Vars:   make(map[string]string),
		Runs:   []actions.Action{},
	}
}

func ParseDirective(s string) (DirectiveType, error) {
	switch s {
	case "run":
		return DIR_RUN, nil
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
