package parser

import "fmt"

type DirectiveType int
type ActionType int

const (
	DIR_RUN DirectiveType = iota
	DIR_CONFIG
	DIR_VAR
)

const (
	ACT_SEED ActionType = iota
	ACT_MIGRATE
	ACT_BACKUP
	ACT_DUMP
)

type Slopfile struct {
	Config map[string]string // config::db.name["x"] → Config["db.name"] = "x"
	Vars   map[string]string // var::seed.rbac["x"] → Vars["seed.rbac"] = "x"
	Runs   []RunAction       // run::seed["x"] → append to Runs
}

type RunAction struct {
	Action ActionType
	Args   string
	Line   int
}

type Directive struct {
	Directive DirectiveType
	Args      string
}

func NewSlopfile() *Slopfile {
	return &Slopfile{
		Config: make(map[string]string),
		Vars:   make(map[string]string),
		Runs:   []RunAction{},
	}
}

func ParseDirective(s string) (DirectiveType, error) {
	switch s {
	case "run":
		return DIR_RUN, nil
	case "config":
		return DIR_CONFIG, nil
	case "var":
		return DIR_VAR, nil
	default:
		return -1, fmt.Errorf("unknown directive: %s", s)
	}
}

func ParseAction(s string) (ActionType, error) {
	switch s {
	case "seed":
		return ACT_SEED, nil
	case "migrate":
		return ACT_MIGRATE, nil
	case "backup":
		return ACT_BACKUP, nil
	case "dump":
		return ACT_DUMP, nil
	default:
		return -1, fmt.Errorf("unknown action: %s", s)
	}
}
