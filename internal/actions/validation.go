package actions

import (
	"fmt"
	"os"
	"strings"
)

var RequiredConfig = map[ActionType][]string{
	ACT_SEED:    {CFG_DB_USER, CFG_DB_NAME, CFG_DB_PASSWORD},
	ACT_MIGRATE: {CFG_DB_USER, CFG_DB_NAME, CFG_DB_PASSWORD},
	ACT_BACKUP:  {CFG_DB_USER, CFG_DB_NAME, CFG_DB_PASSWORD},
	ACT_DUMP:    {CFG_DB_USER, CFG_DB_NAME},
}

var RequiredEnv = map[ActionType][]string{
	ACT_SEED: {"DB_USER", "DB_NAME", "DB_PASSWORD"},
}

func ValidateConfig(action ActionType, config map[string]string) error {
	required, ok := RequiredConfig[action]
	if !ok {
		return nil
	}

	var missing []string
	for _, key := range required {
		if config[key] == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("action %s: missing required config value: %s", action.String(), strings.Join(missing, ", "))
	}

	return nil
}
func ValidateEnv(action ActionType) error {
	required, ok := RequiredEnv[action]
	if !ok {
		return nil
	}

	var missing []string
	for _, envKey := range required {
		if os.Getenv(envKey) == "" {
			missing = append(missing, envKey)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("action %s: missing env vars: %s", action, strings.Join(missing, ", "))
	}
	return nil
}
