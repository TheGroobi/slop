package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/thegroobi/slop/internal/actions"
	"github.com/thegroobi/slop/internal/lexer"
	"github.com/thegroobi/slop/internal/parser"
)

const SLOPFILE_NAME = "Slopfile"

func main() {
	file, err := os.Open(SLOPFILE_NAME)
	if err != nil {
		fmt.Println("Could not read Slopfile", err)
		return
	}
	defer file.Close()

	godotenv.Load()

	fmt.Println("✔ Slopfile loaded 🤤")

	l := lexer.NewLexer(file)
	fmt.Println("✔ Slop lexer initialized!")

	var tokens []lexer.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == lexer.TOKEN_EOF {
			break
		}
	}

	parser := parser.NewParser(tokens)
	fmt.Println("✔ Slop parser initialized!")

	slop, err := parser.Parse()
	if err != nil {
		fmt.Println("Parser failed:", err)
		return
	}

	for len(slop.Runs) > 0 {
		run := slop.Runs[0]
		slop.Runs = slop.Runs[1:]

		var err error
		if len(slop.Config) > 0 {
			err = actions.ValidateConfig(run.Action, slop.Config)
		} else {
			err = actions.ValidateEnv(run.Action)
		}

		if err != nil {
			fmt.Println("Validation error:", err)
			return
		}

		switch run.Action {
		case actions.ACT_TASK:
			t := slop.Tasks[run.Args]
			slop.Runs = append(t, slop.Runs...)
		case actions.ACT_SEED:
			// validate args later (check if valid dir / points to an sql file)
			sa := actions.NewSeedAction(
				run.Args,
				getConfigOrEnv(slop.Config, "db.user", "DB_USER"),
				getConfigOrEnv(slop.Config, "db.name", "DB_NAME"),
				getConfigOrEnv(slop.Config, "db.password", "DB_PASSWORD"),
			)

			if err := run.RunSeed(sa); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("✔ Database %s Seeded properly!\n", slop.Config[actions.CFG_DB_NAME])
		default:
			fmt.Println("Action not valid - might be not implemented yet or missing")
		}
	}
}

func init() {
	godotenv.Load()
}

func getConfigOrEnv(cfg map[string]string, cfgKey, envKey string) string {
	if v := cfg[cfgKey]; v != "" {
		return v
	}
	return os.Getenv(envKey)
}
