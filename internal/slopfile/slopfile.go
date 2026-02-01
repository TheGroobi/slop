package slopfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/thegroobi/slop/internal/actions"
)

type Slopfile struct {
	Config        map[string]string           // config::db.name["x"] → Config["db.name"] = "x"
	Vars          map[string]string           // var::seed.rbac["x"] → Vars["seed.rbac"] = "x"
	Runs          []actions.Action            // run::seed["x"] → append to Runs
	Tasks         map[string][]actions.Action // @task {...} → append to Tasks
	DirectivesRun int
}

func NewSlopfile() *Slopfile {
	return &Slopfile{
		Config:        make(map[string]string),
		Vars:          make(map[string]string),
		Runs:          []actions.Action{},
		Tasks:         make(map[string][]actions.Action),
		DirectivesRun: 0,
	}
}

func (slop *Slopfile) Run(args []string) error {
	err := slop.validateArgs(args)
	if err != nil {
		return err
	}

	if len(args) > 0 {
		fmt.Printf("𝒊 Running %s...\n", args[0])
		err = slop.RunQueue(slop.Tasks[args[0]])
		if err == nil {
			fmt.Printf("✔ Successfully ran %s\n", args[0])
		}
	} else {
		err = slop.RunQueue(slop.Runs)
	}

	return err
}

func (slop *Slopfile) RunQueue(q []actions.Action) error {
	for len(q) > 0 {
		run := q[0]
		q = q[1:]

		var err error
		if len(slop.Config) > 0 {
			err = actions.ValidateConfig(run.Action, slop.Config)
		} else {
			err = actions.ValidateEnv(run.Action)
		}

		if err != nil {
			return err
		}

		switch run.Action {
		case actions.ACT_TASK:
			t := slop.Tasks[run.Args]
			q = append(t, slop.Runs...)
		case actions.ACT_SEED:
			sa := actions.NewSeedAction(
				run.Args,
				getConfigOrEnv(slop.Config, "db.user", "DB_USER"),
				getConfigOrEnv(slop.Config, "db.name", "DB_NAME"),
				getConfigOrEnv(slop.Config, "db.password", "DB_PASSWORD"),
			)

			if err := run.RunSeed(sa); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Action not valid - might be not implemented yet or missing")
		}
		slop.DirectivesRun++
	}

	return nil
}

func getConfigOrEnv(cfg map[string]string, cfgKey, envKey string) string {
	if v := cfg[cfgKey]; v != "" {
		return v
	}
	return os.Getenv(envKey)
}

func (slop *Slopfile) validateArgs(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("Argument invalid: slop only accepts one argument: slop [task]")
	}

	for _, a := range args {
		if _, ok := slop.Tasks[a]; !ok {
			keys := make([]string, 0, len(slop.Tasks))
			for k := range slop.Tasks {
				keys = append(keys, k)
			}

			return fmt.Errorf("Argument invalid: task undefined allowed tasks are: '%s'", strings.Join(keys, ", "))
		}
	}

	return nil
}
