package actions

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ActionType int

const (
	ACT_SEED ActionType = iota
	ACT_MIGRATE
	ACT_BACKUP
	ACT_DUMP
	ACT_ENV

	CFG_DB_USER     = "db.user"
	CFG_DB_NAME     = "db.name"
	CFG_DB_PASSWORD = "db.password"
	CFG_DB_HOST     = "db.host"
	CFG_DB_PORT     = "db.port"
)

type Action struct {
	Action ActionType
	Args   string
	Line   int
}

type SeedAction struct {
	seedDir    string
	dbUser     string
	dbName     string
	dbPassword string
}

func NewSeedAction(seedDir, dbUser, dbName, dbPassword string) *SeedAction {
	return &SeedAction{
		seedDir:    seedDir,
		dbUser:     dbUser,
		dbName:     dbName,
		dbPassword: dbPassword,
	}
}

func (a *Action) RunSeed(s *SeedAction) error {
	var missing []string

	if s.seedDir == "" {
		missing = append(missing, "seed directory")
	}
	if s.dbUser == "" {
		missing = append(missing, "database user")
	}
	if s.dbName == "" {
		missing = append(missing, "database name")
	}
	if s.dbPassword == "" {
		missing = append(missing, "database password")
	}

	if !strings.Contains(s.seedDir, ".sql") {
		return fmt.Errorf("seed: line %d: failed file needs to be a valid SQL file", a.Line)
	}

	if len(missing) > 0 {
		return fmt.Errorf("seed: line %d: missing required fields: %s", a.Line, strings.Join(missing, ", "))
	}

	return runSeedCmd(s)
}

func runSeedCmd(s *SeedAction) error {
	seedCmd := fmt.Sprintf("cat %s | mariadb -u %s -p%s %s", s.seedDir, s.dbUser, s.dbPassword, s.dbName)
	cmd := exec.Command("bash", "-c", seedCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	case "env":
		return ACT_ENV, nil
	default:
		return -1, fmt.Errorf("unknown action: %s", s)
	}
}

func (a ActionType) String() string {
	switch a {
	case ACT_SEED:
		return "seed"
	case ACT_MIGRATE:
		return "migrate"
	case ACT_BACKUP:
		return "backup"
	case ACT_DUMP:
		return "dump"
	default:
		return "unknown"
	}
}
